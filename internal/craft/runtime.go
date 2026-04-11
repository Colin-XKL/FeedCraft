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
	if native, ok, err := buildNativeProcessor(atom, feedURL); err != nil {
		return nil, err
	} else if ok {
		return native, nil
	}

	return buildLegacyProcessor(atom, feedURL)
}

func buildNativeProcessor(atom ResolvedCraftAtom, feedURL string) (engine.FeedProcessor, bool, error) {
	switch atom.TemplateName {
	case "proxy":
		return nil, true, nil
	case "limit":
		maxItems := defaultLimit
		if raw := strings.TrimSpace(atom.Params["num"]); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed <= 0 {
				return nil, false, fmt.Errorf("invalid limit num %q", raw)
			}
			maxItems = parsed
		}
		return &LimitProcessor{MaxItems: maxItems}, true, nil
	case "time-limit":
		days := 7
		if raw := strings.TrimSpace(atom.Params["days"]); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed < 0 {
				return nil, false, fmt.Errorf("invalid time-limit days %q", raw)
			}
			days = parsed
		}
		return &TimeLimitProcessor{Days: days}, true, nil
	case "keyword":
		var mode KeywordFilterMode
		switch strings.TrimSpace(atom.Params["mode"]) {
		case "", string(KeywordIncludeMode):
			mode = KeywordIncludeMode
		case string(KeywordExcludeMode):
			mode = KeywordExcludeMode
		default:
			return nil, false, fmt.Errorf("invalid keyword mode %q", atom.Params["mode"])
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
			return nil, false, fmt.Errorf("invalid keyword scope %q", atom.Params["scope"])
		}

		keywords := splitKeywords(atom.Params["keywords"])
		return &KeywordProcessor{
			Mode:       mode,
			MatchScope: scope,
			Keywords:   keywords,
		}, true, nil
	case "guid-fix":
		return &GUIDFixProcessor{}, true, nil
	case "relative-link-fix":
		return &RelativeLinkFixProcessor{OriginalFeedURL: feedURL}, true, nil
	case "cleanup":
		return newCleanupProcessor(), true, nil
	case "fulltext":
		return newFulltextProcessor(feedURL), true, nil
	case "fulltext-plus":
		return newFulltextPlusProcessor(feedURL, parseFulltextPlusConfig(atom.Params)), true, nil
	case "summary":
		return newSummaryProcessor(atom.Params["prompt"]), true, nil
	case "introduction":
		return newIntroductionProcessor(atom.Params["prompt"]), true, nil
	case "translate-title":
		return newTranslateTitleProcessor(atom.Params["prompt"]), true, nil
	case "translate-content":
		return newTranslateContentProcessor(atom.Params["prompt"]), true, nil
	case "translate-content-immersive":
		return newTranslateContentImmersiveProcessor(atom.Params["prompt"]), true, nil
	case "beautify-content":
		return newBeautifyContentProcessor(atom.Params["prompt"]), true, nil
	case "llm-filter":
		return newLLMFilterProcessor(atom.Params["filter_condition"]), true, nil
	case "ignore-advertorial":
		return newIgnoreAdvertorialProcessor(atom.Params["prompt-for-exclude"]), true, nil
	default:
		return nil, false, nil
	}
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
