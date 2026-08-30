package craft

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"FeedCraft/internal/constant"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"

	"github.com/gorilla/feeds"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"github.com/sirupsen/logrus"
)

type aiFilterResult string

const (
	aiFilterResultKeep aiFilterResult = "keep"
	aiFilterResultDrop aiFilterResult = "drop"
)

type aiFilterDecision struct {
	Reason string         `json:"reason"`
	Result aiFilterResult `json:"result"`
}

type aiFilterExtraPayloadType string

const (
	aiFilterExtraPayloadArticleSummary aiFilterExtraPayloadType = "article_summary"
	aiFilterExtraPayloadArticleContent aiFilterExtraPayloadType = "article_content"
	aiFilterExtraPayloadArticleDate    aiFilterExtraPayloadType = "article_date"
	aiFilterExtraPayloadRawRSSItem     aiFilterExtraPayloadType = "raw_rss_item"
)

var aiFilterCraftParamTmpl = []ParamTemplate{
	{
		Key:         "rule",
		Description: "Rule for deciding which articles should be kept. Example: 只保留科技有关的文章",
		Default:     "",
	},
	{
		Key:         "extra-payload",
		Description: "Comma-separated extra payload list. Supported: article_summary, article_content, article_date, raw_rss_item",
		Default:     string(aiFilterExtraPayloadArticleSummary),
	},
}

func aiFilterCraftLoadParam(m map[string]string) []LegacyCraftOption {
	return GetAIFilterCraftOptions(m["rule"], m["extra-payload"])
}

func GetAIFilterCraftOptions(rule string, extraPayloadRaw string) []LegacyCraftOption {
	return []LegacyCraftOption{
		OptionAIFilter(rule, extraPayloadRaw),
	}
}

func OptionAIFilter(rule string, extraPayloadRaw string) LegacyCraftOption {
	processor := newAIFilterProcessor(rule, extraPayloadRaw)
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		_ = payload
		return applyLocalProcessorToLegacyFeed(context.Background(), processor, feed)
	}
}

type AIFilterProcessor struct {
	rule         string
	payloadTypes []aiFilterExtraPayloadType
}

func newAIFilterProcessor(rule string, extraPayloadRaw string) *AIFilterProcessor {
	return &AIFilterProcessor{
		rule:         strings.TrimSpace(rule),
		payloadTypes: parseAIFilterExtraPayload(extraPayloadRaw),
	}
}

func (p *AIFilterProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	_ = ctx
	if feed == nil || len(feed.Articles) == 0 {
		return feed, nil
	}
	if p.rule == "" {
		return nil, fmt.Errorf("ai-filter requires rule param")
	}

	cloned := cloneCraftFeed(feed)
	drops := parallel.Map(cloned.Articles, func(article *model.CraftArticle, _ int) bool {
		if article == nil {
			return true
		}
		decision, err := p.evaluateAIFilterArticle(article)
		if err != nil {
			logrus.Warnf("failed to evaluate ai-filter for article [%s], err: %v", article.Title, err)
			return false
		}
		return decision.Result == aiFilterResultDrop
	})

	filtered := make([]*model.CraftArticle, 0, len(cloned.Articles))
	for idx, article := range cloned.Articles {
		if article == nil {
			continue
		}
		if !drops[idx] {
			filtered = append(filtered, article)
		}
	}
	cloned.Articles = filtered
	return cloned, nil
}

func (p *AIFilterProcessor) evaluateAIFilterArticle(article *model.CraftArticle) (aiFilterDecision, error) {
	summary := ""
	if lo.Contains(p.payloadTypes, aiFilterExtraPayloadArticleSummary) {
		generated, err := generateAIFilterArticleSummary(article)
		if err != nil {
			return aiFilterDecision{}, err
		}
		summary = generated
	}

	payload, err := buildAIFilterArticlePayload(article, p.payloadTypes, summary)
	if err != nil {
		return aiFilterDecision{}, err
	}

	return cachedAIFilterDecision(article.Title, buildAIFilterPrompt(p.rule), payload)
}

func evaluateAIFilterItem(item *feeds.Item, rule string, payloadTypes []aiFilterExtraPayloadType) (aiFilterDecision, error) {
	processor := &AIFilterProcessor{rule: strings.TrimSpace(rule), payloadTypes: payloadTypes}
	return processor.evaluateAIFilterArticle(articleFromFeedItem(item))
}

func parseAIFilterExtraPayload(raw string) []aiFilterExtraPayloadType {
	return parseAIFilterExtraPayloadWithDefault(raw, []aiFilterExtraPayloadType{aiFilterExtraPayloadArticleSummary})
}

func parseAIFilterExtraPayloadWithDefault(raw string, defaultPayloadTypes []aiFilterExtraPayloadType) []aiFilterExtraPayloadType {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultPayloadTypes
	}

	normalized := strings.NewReplacer("|", ",", "\n", ",", "\t", ",").Replace(raw)
	parts := strings.Split(normalized, ",")
	seen := map[aiFilterExtraPayloadType]bool{}
	payloadTypes := make([]aiFilterExtraPayloadType, 0, len(parts))
	for _, part := range parts {
		payloadType := aiFilterExtraPayloadType(strings.TrimSpace(part))
		switch payloadType {
		case aiFilterExtraPayloadArticleSummary, aiFilterExtraPayloadArticleContent, aiFilterExtraPayloadArticleDate, aiFilterExtraPayloadRawRSSItem:
			if !seen[payloadType] {
				seen[payloadType] = true
				payloadTypes = append(payloadTypes, payloadType)
			}
		case "":
			continue
		default:
			logrus.Warnf("unknown ai-filter extra-payload value [%s], ignored", payloadType)
		}
	}
	if len(payloadTypes) == 0 {
		return defaultPayloadTypes
	}
	return payloadTypes
}

func buildAIFilterPrompt(rule string) string {
	return fmt.Sprintf(`You are an RSS article filtering assistant.

Rule from user:
%s

Decide whether the article should be kept in the RSS feed or dropped from the RSS feed.

Output requirements:
- Return JSON only. Do not include markdown fences or any other text.
- JSON schema: {"reason":"short reason","result":"keep|drop"}
- Use result="keep" when the article should remain visible to the user.
- Use result="drop" when the article should be excluded from the RSS output.

Examples:
{"reason":"The article is about semiconductor technology and matches the rule.","result":"keep"}
{"reason":"The article is unrelated to the requested topic.","result":"drop"}`, rule)
}

func cachedAIFilterDecision(title string, prompt string, context string) (aiFilterDecision, error) {
	hashVal := util.GetTextContentHash(strings.Join([]string{
		util.GetTextContentHash(prompt),
		util.GetTextContentHash(context),
	}, "|"))
	cacheKey := getCraftCacheKey("ai-filter", hashVal)

	cached, err := util.CachedFuncWithPreLog(cacheKey, func() (string, error) {
		result, err := llmContextCaller(prompt, context, util.ContentProcessOption{
			RemoveImage: true,
			ConvertToMd: true,
			Temperature: util.LowestLLMTemperaturePtr(),
		})
		if err != nil {
			return "", err
		}
		decision, err := parseAIFilterDecision(result)
		if err != nil {
			return "", err
		}
		normalized, err := json.Marshal(decision)
		if err != nil {
			return "", err
		}
		return string(normalized), nil
	}, func(isCached bool) {
		logrus.Infof("applying craft [ai-filter] to article [%s], cached: %v", title, isCached)
	})
	if err != nil {
		return aiFilterDecision{}, err
	}
	return parseAIFilterDecision(cached)
}

func buildAIFilterArticlePayload(article *model.CraftArticle, payloadTypes []aiFilterExtraPayloadType, articleSummary string) (string, error) {
	if article == nil {
		return "", fmt.Errorf("nil rss item")
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "Article Title:\n```text\n%s\n```", strings.TrimSpace(article.Title))

	for _, payloadType := range payloadTypes {
		switch payloadType {
		case aiFilterExtraPayloadArticleSummary:
			if strings.TrimSpace(articleSummary) != "" {
				fmt.Fprintf(&builder, "\n\nArticle Summary:\n```markdown\n%s\n```", strings.TrimSpace(articleSummary))
			}
		case aiFilterExtraPayloadArticleContent:
			fmt.Fprintf(&builder, "\n\nArticle Content:\n```markdown\n%s\n```", strings.TrimSpace(getPrimaryArticleContent(article)))
		case aiFilterExtraPayloadArticleDate:
			builder.WriteString("\n\nArticle Date:\n```text\n")
			if !article.Created.IsZero() {
				fmt.Fprintf(&builder, "Created: %s\n", article.Created.Format("2006-01-02T15:04:05Z07:00"))
			}
			if !article.Updated.IsZero() {
				fmt.Fprintf(&builder, "Updated: %s\n", article.Updated.Format("2006-01-02T15:04:05Z07:00"))
			}
			builder.WriteString("```")
		case aiFilterExtraPayloadRawRSSItem:
			rawJSON, err := buildRawRSSItemJSON(article)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&builder, "\n\nRaw RSS Item JSON:\n```json\n%s\n```", rawJSON)
		}
	}

	return builder.String(), nil
}

func generateAIFilterArticleSummary(article *model.CraftArticle) (string, error) {
	content := getPrimaryArticleContent(article)
	if strings.TrimSpace(content) == "" {
		return "", nil
	}

	summaryPrompt := renderTargetLangPrompt("", constant.DefaultPrompts[constant.ProcessorTypeSummary])
	hashVal := util.GetTextContentHash(strings.Join([]string{
		util.GetTextContentHash(summaryPrompt),
		strings.TrimSpace(article.Title),
		util.GetTextContentHash(content),
	}, "|"))
	cacheKey := getCraftCacheKey("ai-filter-article-summary", hashVal)

	return util.CachedFuncWithPreLog(cacheKey, func() (string, error) {
		processedContent := content
		domain, _ := util.ParseDomainFromUrl(article.Link)
		// Drop inline base64 images before converting: they carry no useful signal
		// for LLM processing but can dominate the token budget.
		cleanedContent := util.HTMLToMarkdown(util.RemoveBase64Images(content), domain)
		if strings.TrimSpace(cleanedContent) != "" {
			processedContent = cleanedContent
		}
		return llmContextCaller(summaryPrompt, processedContent, util.ContentProcessOption{
			Temperature: util.LowestLLMTemperaturePtr(),
		})
	}, func(isCached bool) {
		logrus.Infof("generating ai-filter article summary for article [%s], cached: %v", article.Title, isCached)
	})
}

func parseAIFilterDecision(raw string) (aiFilterDecision, error) {
	jsonText, err := extractAIFilterJSON(raw)
	if err != nil {
		return aiFilterDecision{}, err
	}

	var decision aiFilterDecision
	if err := json.Unmarshal([]byte(jsonText), &decision); err != nil {
		return aiFilterDecision{}, err
	}
	decision.Result = aiFilterResult(strings.ToLower(strings.TrimSpace(string(decision.Result))))
	decision.Reason = strings.TrimSpace(decision.Reason)
	switch decision.Result {
	case aiFilterResultKeep, aiFilterResultDrop:
		return decision, nil
	default:
		return aiFilterDecision{}, fmt.Errorf("unexpected ai-filter result [%s]", decision.Result)
	}
}

func extractAIFilterJSON(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "```") {
		lines := strings.Split(trimmed, "\n")
		if len(lines) >= 3 {
			trimmed = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start < 0 || end < start {
		return "", fmt.Errorf("ai-filter response does not contain json object: [%s]", raw)
	}
	return strings.TrimSpace(trimmed[start : end+1]), nil
}

func buildRawRSSItemJSON(article *model.CraftArticle) (string, error) {
	raw := map[string]string{
		"title":       article.Title,
		"description": article.Description,
		"content":     article.Content,
		"id":          article.Id,
	}
	if article.Link != "" {
		raw["link"] = article.Link
	}
	if article.Source != "" {
		raw["source"] = article.Source
	}
	if article.AuthorName != "" || article.AuthorEmail != "" {
		raw["author_name"] = article.AuthorName
		raw["author_email"] = article.AuthorEmail
	}
	if !article.Created.IsZero() {
		raw["created"] = article.Created.Format("2006-01-02T15:04:05Z07:00")
	}
	if !article.Updated.IsZero() {
		raw["updated"] = article.Updated.Format("2006-01-02T15:04:05Z07:00")
	}

	encoded, err := json.Marshal(raw)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
