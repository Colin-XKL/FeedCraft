package craft

import (
	"context"
	"fmt"
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

type nativeProcessorBuilder func(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error)

var nativeProcessorBuilders = map[string]nativeProcessorBuilder{
	"proxy":                       buildNativeProxyProcessor,
	"limit":                       buildNativeLimitProcessor,
	"time-limit":                  buildNativeTimeLimitProcessor,
	"keyword":                     buildNativeKeywordProcessor,
	"guid-fix":                    buildNativeGUIDFixProcessor,
	"relative-link-fix":           buildNativeRelativeLinkFixProcessor,
	"cleanup":                     buildNativeCleanupProcessor,
	"fulltext":                    buildNativeFulltextProcessor,
	"fulltext-plus":               buildNativeFulltextPlusProcessor,
	"summary":                     buildNativeSummaryProcessor,
	"introduction":                buildNativeIntroductionProcessor,
	"translate-title":             buildNativeTranslateTitleProcessor,
	"translate-content":           buildNativeTranslateContentProcessor,
	"translate-content-immersive": buildNativeTranslateContentImmersiveProcessor,
	"beautify-content":            buildNativeBeautifyContentProcessor,
	"llm-filter":                  buildNativeLLMFilterProcessor,
	"ignore-advertorial":          buildNativeIgnoreAdvertorialProcessor,
}

type NoopProcessor struct{}

func (p *NoopProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	return feed, nil
}

type LimitProcessor struct {
	MaxItems int
}

func (p *LimitProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || p == nil || p.MaxItems <= 0 || len(feed.Articles) <= p.MaxItems {
		return feed, nil
	}
	cloned := cloneCraftFeed(feed)
	cloned.Articles = cloned.Articles[:p.MaxItems]
	return cloned, nil
}

type TimeLimitProcessor struct {
	Days int
	Now  func() time.Time
}

func (p *TimeLimitProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || p == nil {
		return feed, nil
	}
	nowFunc := p.Now
	if nowFunc == nil {
		nowFunc = time.Now
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

	cutoff := nowFunc().AddDate(0, 0, -p.Days)
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
}

type KeywordProcessor struct {
	Mode       KeywordFilterMode
	MatchScope KeywordMatchScope
	Keywords   []string
}

func (p *KeywordProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || p == nil || len(p.Keywords) == 0 {
		return feed, nil
	}

	searchTitle := p.MatchScope == KeywordMatchAll || p.MatchScope == KeywordMatchTitle
	searchContent := p.MatchScope == KeywordMatchAll || p.MatchScope == KeywordMatchContent
	cloned := cloneCraftFeed(feed)
	filtered := make([]*model.CraftArticle, 0, len(cloned.Articles))

	for _, article := range cloned.Articles {
		if article == nil {
			continue
		}
		matched := false
		for _, keyword := range p.Keywords {
			if searchTitle && strings.Contains(article.Title, keyword) {
				matched = true
				break
			}
			if searchContent && (strings.Contains(article.Content, keyword) || strings.Contains(article.Description, keyword)) {
				matched = true
				break
			}
		}

		switch p.Mode {
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
}

type GUIDFixProcessor struct{}

func (p *GUIDFixProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
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
		article.Id = fmt.Sprintf("%x", util.GetMD5Hash(article.Title+article.Content+article.Description))
	}
	return cloned, nil
}

type RelativeLinkFixProcessor struct {
	OriginalFeedURL string
}

func (p *RelativeLinkFixProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || p == nil {
		return feed, nil
	}

	cloned := cloneCraftFeed(feed)
	for _, article := range cloned.Articles {
		if article == nil || strings.TrimSpace(article.Link) == "" {
			continue
		}
		absURL := getAbsLinkForFeedItem(p.OriginalFeedURL, cloned.Link, article.Link)
		if absURL != "" {
			article.Link = absURL
		}
	}
	return cloned, nil
}

func BuildProcessor(db *gorm.DB, craftName string, feedURL string) (engine.FeedProcessor, error) {
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

	processors := make([]engine.FeedProcessor, 0, len(atoms))
	for _, atom := range atoms {
		processor, err := buildProcessorForAtom(atom, feedURL)
		if err != nil {
			return nil, err
		}
		if processor != nil {
			processors = append(processors, processor)
		}
	}
	if len(processors) == 0 {
		return nil, nil
	}
	return &engine.FlowCraftProcessor{Processors: processors}, nil
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

func buildProcessorForAtom(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	builder, ok := nativeProcessorBuilders[atom.TemplateName]
	if ok {
		return builder(atom, feedURL)
	}

	return buildLegacyProcessor(atom, feedURL)
}

func buildNativeProxyProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return &NoopProcessor{}, nil
}

func buildNativeLimitProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	maxItems := defaultLimit
	if raw := strings.TrimSpace(atom.Params["num"]); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			return nil, fmt.Errorf("invalid limit num %q", raw)
		}
		maxItems = parsed
	}
	return &LimitProcessor{MaxItems: maxItems}, nil
}

func buildNativeTimeLimitProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	days := 7
	if raw := strings.TrimSpace(atom.Params["days"]); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 0 {
			return nil, fmt.Errorf("invalid time-limit days %q", raw)
		}
		days = parsed
	}
	return &TimeLimitProcessor{Days: days}, nil
}

func buildNativeKeywordProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
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

	return &KeywordProcessor{
		Mode:       mode,
		MatchScope: scope,
		Keywords:   splitKeywords(atom.Params["keywords"]),
	}, nil
}

func buildNativeGUIDFixProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return &GUIDFixProcessor{}, nil
}

func buildNativeRelativeLinkFixProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return &RelativeLinkFixProcessor{OriginalFeedURL: feedURL}, nil
}

func buildNativeCleanupProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newCleanupProcessor(), nil
}

func buildNativeFulltextProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newFulltextProcessor(feedURL), nil
}

func buildNativeFulltextPlusProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newFulltextPlusProcessor(feedURL, parseFulltextPlusConfig(atom.Params)), nil
}

func buildNativeSummaryProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newSummaryProcessor(atom.Params["prompt"]), nil
}

func buildNativeIntroductionProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newIntroductionProcessor(atom.Params["prompt"]), nil
}

func buildNativeTranslateTitleProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newTranslateTitleProcessor(atom.Params["prompt"]), nil
}

func buildNativeTranslateContentProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newTranslateContentProcessor(atom.Params["prompt"]), nil
}

func buildNativeTranslateContentImmersiveProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newTranslateContentImmersiveProcessor(atom.Params["prompt"]), nil
}

func buildNativeBeautifyContentProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newBeautifyContentProcessor(atom.Params["prompt"]), nil
}

func buildNativeLLMFilterProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newLLMFilterProcessor(atom.Params["filter_condition"]), nil
}

func buildNativeIgnoreAdvertorialProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
	return newIgnoreAdvertorialProcessor(atom.Params["prompt-for-exclude"]), nil
}

func buildLegacyProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, error) {
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
	processors := make([]engine.FeedProcessor, 0, len(options))
	for _, opt := range options {
		processors = append(processors, &LegacyOptionAdapter{
			Option: opt,
			Extra:  payload,
		})
	}
	return &engine.FlowCraftProcessor{Processors: processors}, nil
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
