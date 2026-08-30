package craft

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/engine"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"

	"gorm.io/gorm"
)

type ResolvedCraftAtom struct {
	Name         string
	TemplateName string
	Params       map[string]string
}

type nativeProcessorBuilder func(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error)

func buildZeroArg(ctor func() localProcessor) nativeProcessorBuilder {
	return func(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
		_ = atom
		_ = feedURL
		return wrapLocalProcessor(ctor()), nil
	}
}

func buildWithStringParam(paramKey string, ctor func(string) localProcessor) nativeProcessorBuilder {
	return func(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
		_ = feedURL
		return wrapLocalProcessor(ctor(atom.Params[paramKey])), nil
	}
}

func buildWithFeedURL(ctor func(string) localProcessor) nativeProcessorBuilder {
	return func(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
		_ = atom
		return wrapLocalProcessor(ctor(feedURL)), nil
	}
}

var nativeProcessorBuilders = map[string]nativeProcessorBuilder{
	"proxy":                       buildNativeProxyProcessor,
	"limit":                       buildNativeLimitProcessor,
	"time-limit":                  buildNativeTimeLimitProcessor,
	"keyword":                     buildNativeKeywordProcessor,
	"guid-fix":                    buildNativeGUIDFixProcessor,
	"relative-link-fix":           buildNativeRelativeLinkFixProcessor,
	"link-flatten":                buildNativeLinkFlattenProcessor,
	"cleanup":                     buildZeroArg(func() localProcessor { return newCleanupProcessor() }),
	"fulltext":                    buildWithFeedURL(func(u string) localProcessor { return newFulltextProcessor(u) }),
	"fulltext-plus":               buildNativeFulltextPlusProcessor,
	"summary":                     buildWithStringParam("prompt", func(s string) localProcessor { return newSummaryProcessor(s) }),
	"introduction":                buildWithStringParam("prompt", func(s string) localProcessor { return newIntroductionProcessor(s) }),
	"re-title":                    buildWithStringParam("prompt", func(s string) localProcessor { return newRetitleProcessor(s) }),
	"translate-title":             buildWithStringParam("prompt", func(s string) localProcessor { return newTranslateTitleProcessor(s) }),
	"translate-content":           buildWithStringParam("prompt", func(s string) localProcessor { return newTranslateContentProcessor(s) }),
	"translate-content-immersive": buildWithStringParam("prompt", func(s string) localProcessor { return newTranslateContentImmersiveProcessor(s) }),
	"beautify-content":            buildWithStringParam("prompt", func(s string) localProcessor { return newBeautifyContentProcessor(s) }),
	"llm-filter":                  buildWithStringParam("filter_condition", func(s string) localProcessor { return newLLMFilterProcessor(s) }),
	"ignore-advertorial":          buildWithStringParam("prompt-for-exclude", func(s string) localProcessor { return newIgnoreAdvertorialProcessor(s) }),
	"embedding-filter":            buildNativeEmbeddingFilterProcessor,
	"ai-filter":                   buildNativeAIFilterProcessor,
	"ai-content-process":          buildNativeAIContentProcessProcessor,
}

func BuildOptionChain(db *gorm.DB, craftName string, feedURL string) (engine.CraftOption, error) {
	if db == nil {
		db = util.GetDatabase()
	}

	atoms, err := ResolveCraftAtoms(db, craftName)
	if err != nil {
		return nil, err
	}
	if len(atoms) == 0 {
		return nil, nil
	}

	options := make([]engine.CraftOption, 0, len(atoms))
	for _, atom := range atoms {
		option, err := buildOptionForAtom(atom, feedURL)
		if err != nil {
			return nil, err
		}
		if option != nil {
			options = append(options, option)
		}
	}
	if len(options) == 0 {
		return nil, nil
	}
	return composeEngineOptions(options...), nil
}

func ResolveCraftAtoms(db *gorm.DB, craftName string) ([]ResolvedCraftAtom, error) {
	if db == nil {
		db = util.GetDatabase()
	}

	craftAtomDict := getCraftAtomDict(db)
	craftTmplDict := GetSysCraftTemplateDict()
	return resolveCraftAtoms(db, &craftAtomDict, &craftTmplDict, craftName, 0)
}

func resolveCraftAtoms(db *gorm.DB, craftAtomDict *map[string]dao.CraftAtom, craftTmplDict *map[string]CraftTemplate, craftName string, depthID int) ([]ResolvedCraftAtom, error) {
	if depthID+1 > MaxCallDepth {
		return nil, fmt.Errorf("max call depth hit")
	}

	if strings.Contains(craftName, ",") {
		parts := strings.Split(craftName, ",")
		resolved := make([]ResolvedCraftAtom, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			sub, err := resolveCraftAtoms(db, craftAtomDict, craftTmplDict, part, depthID)
			if err != nil {
				return nil, err
			}
			resolved = append(resolved, sub...)
		}
		return resolved, nil
	}

	craftAtom, isKnownCraftAtom := (*craftAtomDict)[craftName]
	if isKnownCraftAtom {
		if _, tmplValid := (*craftTmplDict)[craftAtom.TemplateName]; !tmplValid {
			return nil, fmt.Errorf("invalid tmpl name [%s] for craft atom [%s]", craftAtom.TemplateName, craftAtom.Name)
		}
		return []ResolvedCraftAtom{{
			Name:         craftAtom.Name,
			TemplateName: craftAtom.TemplateName,
			Params:       cloneParams(craftAtom.Params),
		}}, nil
	}

	craftArr, err := extractCraftArrFromFlow(db, craftName)
	if err != nil {
		return nil, fmt.Errorf("not a valid craft name")
	}

	resolved := make([]ResolvedCraftAtom, 0, len(craftArr))
	for _, extractedSubCraftName := range craftArr {
		sub, recurErr := resolveCraftAtoms(db, craftAtomDict, craftTmplDict, extractedSubCraftName, depthID+1)
		if recurErr != nil {
			return nil, recurErr
		}
		resolved = append(resolved, sub...)
	}
	return resolved, nil
}

func buildOptionForAtom(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	builder, ok := nativeProcessorBuilders[atom.TemplateName]
	if ok {
		return builder(atom, feedURL)
	}

	return buildLegacyOption(atom, feedURL)
}

func buildNativeProxyProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		_ = ctx
		return feed, nil
	}, nil
}

func craftArticleTime(article *model.CraftArticle) time.Time {
	if article == nil {
		return time.Time{}
	}
	if !article.Created.IsZero() {
		return article.Created
	}
	return article.Updated
}

func buildNativeLimitProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	maxItems := defaultLimit
	if raw := strings.TrimSpace(atom.Params["num"]); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			return nil, fmt.Errorf("invalid limit num %q", raw)
		}
		maxItems = parsed
	}
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		_ = ctx
		if feed == nil || maxItems <= 0 || len(feed.Articles) == 0 {
			return feed, nil
		}
		cloned := cloneCraftFeed(feed)
		sort.SliceStable(cloned.Articles, func(i, j int) bool {
			return craftArticleTime(cloned.Articles[i]).After(craftArticleTime(cloned.Articles[j]))
		})
		if len(cloned.Articles) > maxItems {
			cloned.Articles = cloned.Articles[:maxItems]
		}
		return cloned, nil
	}, nil
}

func buildNativeTimeLimitProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	days := 7
	if raw := strings.TrimSpace(atom.Params["days"]); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 0 {
			return nil, fmt.Errorf("invalid time-limit days %q", raw)
		}
		days = parsed
	}
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		_ = ctx
		if feed == nil {
			return feed, nil
		}
		hasNormalDate := false
		for _, article := range feed.Articles {
			if article == nil {
				continue
			}
			if !article.Created.IsZero() && article.Created.Year() > 1970 {
				hasNormalDate = true
				break
			}
		}
		if !hasNormalDate {
			return feed, nil
		}
		cutoff := time.Now().AddDate(0, 0, -days)
		cloned := cloneCraftFeed(feed)
		filtered := make([]*model.CraftArticle, 0, len(cloned.Articles))
		for _, article := range cloned.Articles {
			if article == nil {
				continue
			}
			if article.Created.IsZero() {
				filtered = append(filtered, article)
				continue
			}
			if article.Created.Year() <= 1970 {
				continue
			}
			if !article.Created.Before(cutoff) {
				filtered = append(filtered, article)
			}
		}
		cloned.Articles = filtered
		return cloned, nil
	}, nil
}

func buildNativeKeywordProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	var mode KeywordFilterMode
	switch strings.TrimSpace(atom.Params["mode"]) {
	case "", string(KeywordIncludeMode):
		mode = KeywordIncludeMode
	case string(KeywordExcludeMode):
		mode = KeywordExcludeMode
	default:
		return nil, fmt.Errorf("invalid keyword mode %q", atom.Params["mode"])
	}

	var scope KeywordMatchScope
	switch strings.TrimSpace(atom.Params["scope"]) {
	case "", string(KeywordMatchAll):
		scope = KeywordMatchAll
	case string(KeywordMatchTitle):
		scope = KeywordMatchTitle
	case string(KeywordMatchContent):
		scope = KeywordMatchContent
	default:
		return nil, fmt.Errorf("invalid keyword scope %q", atom.Params["scope"])
	}
	keywords := splitKeywords(atom.Params["keywords"])
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		_ = ctx
		if feed == nil || len(keywords) == 0 {
			return feed, nil
		}
		searchTitle := scope == KeywordMatchAll || scope == KeywordMatchTitle
		searchContent := scope == KeywordMatchAll || scope == KeywordMatchContent
		cloned := cloneCraftFeed(feed)
		filtered := make([]*model.CraftArticle, 0, len(cloned.Articles))
		for _, article := range cloned.Articles {
			if article == nil {
				continue
			}
			matched := false
			for _, keyword := range keywords {
				if searchTitle && strings.Contains(article.Title, keyword) {
					matched = true
					break
				}
				if searchContent && (strings.Contains(article.Content, keyword) || strings.Contains(article.Description, keyword)) {
					matched = true
					break
				}
			}
			switch mode {
			case KeywordIncludeMode:
				if matched {
					filtered = append(filtered, article)
				}
			case KeywordExcludeMode:
				if !matched {
					filtered = append(filtered, article)
				}
			default:
				filtered = append(filtered, article)
			}
		}
		cloned.Articles = filtered
		return cloned, nil
	}, nil
}

func buildNativeGUIDFixProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		_ = ctx
		if feed == nil {
			return feed, nil
		}
		cloned := cloneCraftFeed(feed)
		for _, article := range cloned.Articles {
			if article == nil {
				continue
			}
			if article.Title == "" && article.Content == "" && article.Description == "" {
				article.Id = article.Link
				continue
			}
			article.Id = util.GetTextContentHash(article.Title + article.Content + article.Description)
		}
		return cloned, nil
	}, nil
}

func buildNativeRelativeLinkFixProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		_ = ctx
		if feed == nil {
			return feed, nil
		}
		cloned := cloneCraftFeed(feed)
		for _, article := range cloned.Articles {
			if article == nil || strings.TrimSpace(article.Link) == "" {
				continue
			}
			absURL := getAbsLinkForFeedItem(feedURL, cloned.Link, article.Link)
			if absURL != "" {
				article.Link = absURL
			}
		}
		return cloned, nil
	}, nil
}

func buildNativeFulltextPlusProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	return wrapLocalProcessor(newFulltextPlusProcessor(feedURL, parseFulltextPlusConfig(atom.Params))), nil
}

func buildNativeEmbeddingFilterProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	_ = feedURL
	cfg, err := parseEmbeddingFilterParams(atom.Params)
	if err != nil {
		return nil, err
	}
	return wrapLocalProcessor(NewEmbeddingFilterProcessor(cfg.anchors, cfg.threshold, cfg.maxContentLen, cfg.instruction, cfg.mode)), nil
}

func buildNativeAIFilterProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	_ = feedURL
	return wrapLocalProcessor(newAIFilterProcessor(atom.Params["rule"], atom.Params["extra-payload"])), nil
}

func buildNativeAIContentProcessProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	_ = feedURL
	return wrapLocalProcessor(newAIContentProcessProcessor(atom.Params["rule"], atom.Params["extra-payload"], atom.Params["placement"])), nil
}

func buildLegacyOption(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	tmplDict := GetSysCraftTemplateDict()
	tmpl, ok := tmplDict[atom.TemplateName]
	if !ok {
		return nil, fmt.Errorf("invalid tmpl name [%s] for craft atom [%s]", atom.TemplateName, atom.Name)
	}

	options := tmpl.GetOptions(atom.Params)
	if len(options) == 0 {
		return nil, nil
	}

	payload := ExtraPayload{originalFeedUrl: feedURL}
	adapted := make([]engine.CraftOption, 0, len(options))
	for _, opt := range options {
		adapted = append(adapted, wrapLocalOption(AdaptLegacyOption(opt, payload)))
	}
	return composeEngineOptions(adapted...), nil
}

func wrapLocalOption(option CraftOption) engine.CraftOption {
	if option == nil {
		return nil
	}
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		return option(ctx, feed)
	}
}

type localProcessor interface {
	Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error)
}

func wrapLocalProcessor(processor localProcessor) engine.CraftOption {
	if processor == nil {
		return nil
	}
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		return processor.Process(ctx, feed)
	}
}

func WrapLocalProcessor(processor localProcessor) engine.CraftOption {
	return wrapLocalProcessor(processor)
}

func composeEngineOptions(options ...engine.CraftOption) engine.CraftOption {
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		currentFeed := feed
		var err error
		for _, option := range options {
			if option == nil {
				continue
			}
			currentFeed, err = option(ctx, currentFeed)
			if err != nil {
				return nil, err
			}
		}
		return currentFeed, nil
	}
}

func cloneCraftFeed(feed *model.CraftFeed) *model.CraftFeed {
	if feed == nil {
		return nil
	}
	cloned := *feed
	cloned.Articles = make([]*model.CraftArticle, 0, len(feed.Articles))
	for _, article := range feed.Articles {
		if article == nil {
			cloned.Articles = append(cloned.Articles, nil)
			continue
		}
		articleCopy := *article
		cloned.Articles = append(cloned.Articles, &articleCopy)
	}
	return &cloned
}

func cloneParams(params map[string]string) map[string]string {
	if len(params) == 0 {
		return map[string]string{}
	}
	cloned := make(map[string]string, len(params))
	for key, value := range params {
		cloned[key] = value
	}
	return cloned
}

func splitKeywords(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	keywords := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			keywords = append(keywords, trimmed)
		}
	}
	return keywords
}
