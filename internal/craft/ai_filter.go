package craft

import (
	"encoding/json"
	"fmt"
	"strings"

	"FeedCraft/internal/constant"
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

func aiFilterCraftLoadParam(m map[string]string) []CraftOption {
	return GetAIFilterCraftOptions(m["rule"], m["extra-payload"])
}

func GetAIFilterCraftOptions(rule string, extraPayloadRaw string) []CraftOption {
	return []CraftOption{
		OptionAIFilter(rule, extraPayloadRaw),
	}
}

func OptionAIFilter(rule string, extraPayloadRaw string) CraftOption {
	rule = strings.TrimSpace(rule)
	payloadTypes := parseAIFilterExtraPayload(extraPayloadRaw)
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		if rule == "" {
			return fmt.Errorf("ai-filter requires rule param")
		}
		items := feed.Items
		if len(items) == 0 {
			return nil
		}

		drops := parallel.Map(items, func(item *feeds.Item, _ int) bool {
			decision, err := evaluateAIFilterItem(item, rule, payloadTypes)
			if err != nil {
				logrus.Warnf("failed to evaluate ai-filter for article [%s], err: %v", item.Title, err)
				return false
			}
			return decision.Result == aiFilterResultDrop
		})

		feed.Items = lo.Filter(items, func(_ *feeds.Item, index int) bool {
			return !drops[index]
		})
		return nil
	}
}

func evaluateAIFilterItem(item *feeds.Item, rule string, payloadTypes []aiFilterExtraPayloadType) (aiFilterDecision, error) {
	summary := ""
	if lo.Contains(payloadTypes, aiFilterExtraPayloadArticleSummary) {
		generated, err := generateAIFilterArticleSummary(item)
		if err != nil {
			return aiFilterDecision{}, err
		}
		summary = generated
	}

	context, err := buildAIFilterArticlePayload(item, payloadTypes, summary)
	if err != nil {
		return aiFilterDecision{}, err
	}

	result, err := llmContextCaller(buildAIFilterPrompt(rule), context, util.ContentProcessOption{
		RemoveImage: true,
		ConvertToMd: true,
	})
	if err != nil {
		return aiFilterDecision{}, err
	}
	return parseAIFilterDecision(result)
}

func parseAIFilterExtraPayload(raw string) []aiFilterExtraPayloadType {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []aiFilterExtraPayloadType{aiFilterExtraPayloadArticleSummary}
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
		return []aiFilterExtraPayloadType{aiFilterExtraPayloadArticleSummary}
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

func buildAIFilterArticlePayload(item *feeds.Item, payloadTypes []aiFilterExtraPayloadType, articleSummary string) (string, error) {
	if item == nil {
		return "", fmt.Errorf("nil rss item")
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "Article Title:\n```text\n%s\n```", strings.TrimSpace(item.Title))

	for _, payloadType := range payloadTypes {
		switch payloadType {
		case aiFilterExtraPayloadArticleSummary:
			if strings.TrimSpace(articleSummary) != "" {
				fmt.Fprintf(&builder, "\n\nArticle Summary:\n```markdown\n%s\n```", strings.TrimSpace(articleSummary))
			}
		case aiFilterExtraPayloadArticleContent:
			fmt.Fprintf(&builder, "\n\nArticle Content:\n```markdown\n%s\n```", strings.TrimSpace(getPrimaryFeedItemContent(item)))
		case aiFilterExtraPayloadArticleDate:
			builder.WriteString("\n\nArticle Date:\n```text\n")
			if !item.Created.IsZero() {
				fmt.Fprintf(&builder, "Created: %s\n", item.Created.Format("2006-01-02T15:04:05Z07:00"))
			}
			if !item.Updated.IsZero() {
				fmt.Fprintf(&builder, "Updated: %s\n", item.Updated.Format("2006-01-02T15:04:05Z07:00"))
			}
			builder.WriteString("```")
		case aiFilterExtraPayloadRawRSSItem:
			rawJSON, err := buildRawRSSItemJSON(item)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&builder, "\n\nRaw RSS Item JSON:\n```json\n%s\n```", rawJSON)
		}
	}

	return builder.String(), nil
}

func generateAIFilterArticleSummary(item *feeds.Item) (string, error) {
	content := getPrimaryFeedItemContent(item)
	if strings.TrimSpace(content) == "" {
		return "", nil
	}

	summaryPrompt := renderTargetLangPrompt("", constant.DefaultPrompts[constant.ProcessorTypeSummary])
	hashVal := util.GetTextContentHash(strings.Join([]string{
		util.GetTextContentHash(summaryPrompt),
		strings.TrimSpace(item.Title),
		util.GetTextContentHash(content),
	}, "|"))
	cacheKey := getCraftCacheKey("ai-filter-article-summary", hashVal)

	return util.CachedFuncWithPreLog(cacheKey, func() (string, error) {
		processedContent := content
		domain := ""
		if item.Link != nil {
			domain, _ = util.ParseDomainFromUrl(item.Link.Href)
		}
		cleanedContent := util.Html2Markdown(content, &domain)
		if strings.TrimSpace(cleanedContent) != "" {
			processedContent = cleanedContent
		}
		return llmContextCaller(summaryPrompt, processedContent, util.ContentProcessOption{})
	}, func(isCached bool) {
		logrus.Infof("generating ai-filter article summary for article [%s], cached: %v", item.Title, isCached)
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

func getPrimaryFeedItemContent(item *feeds.Item) string {
	if item == nil {
		return ""
	}
	content := item.Content
	if strings.TrimSpace(content) == "" {
		content = item.Description
	}
	return content
}

func buildRawRSSItemJSON(item *feeds.Item) (string, error) {
	raw := map[string]string{
		"title":       item.Title,
		"description": item.Description,
		"content":     item.Content,
		"id":          item.Id,
	}
	if item.Link != nil {
		raw["link"] = item.Link.Href
	}
	if item.Source != nil {
		raw["source"] = item.Source.Href
	}
	if item.Author != nil {
		raw["author_name"] = item.Author.Name
		raw["author_email"] = item.Author.Email
	}
	if !item.Created.IsZero() {
		raw["created"] = item.Created.Format("2006-01-02T15:04:05Z07:00")
	}
	if !item.Updated.IsZero() {
		raw["updated"] = item.Updated.Format("2006-01-02T15:04:05Z07:00")
	}

	encoded, err := json.Marshal(raw)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
