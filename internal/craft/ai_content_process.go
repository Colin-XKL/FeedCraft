package craft

import (
	"fmt"
	"strings"

	"FeedCraft/internal/util"

	"github.com/gorilla/feeds"
)

type aiContentProcessPlacement string

const (
	aiContentProcessPlacementPrepend aiContentProcessPlacement = "prepend"
	aiContentProcessPlacementReplace aiContentProcessPlacement = "replace"
	aiContentProcessPlacementAppend  aiContentProcessPlacement = "append"
)

var aiContentProcessCraftParamTmpl = []ParamTemplate{
	{
		Key:         "rule",
		Description: "Instruction for processing each article content. Example: 总结文章的关键观点并列出行动建议",
		Default:     "",
	},
	{
		Key:         "extra-payload",
		Description: "Comma-separated extra payload list. Supported: article_summary, article_content, article_date, raw_rss_item",
		Default:     string(aiFilterExtraPayloadArticleContent),
	},
	{
		Key:         "placement",
		Description: "Where to write generated content. Supported: prepend, replace, append",
		Default:     string(aiContentProcessPlacementPrepend),
	},
}

func aiContentProcessCraftLoadParam(m map[string]string) []CraftOption {
	return GetAIContentProcessCraftOptions(m["rule"], m["extra-payload"], m["placement"])
}

func GetAIContentProcessCraftOptions(rule string, extraPayloadRaw string, placementRaw string) []CraftOption {
	return []CraftOption{
		OptionAIContentProcess(rule, extraPayloadRaw, placementRaw),
	}
}

func OptionAIContentProcess(rule string, extraPayloadRaw string, placementRaw string) CraftOption {
	rule = strings.TrimSpace(rule)
	payloadTypes := parseAIFilterExtraPayloadWithDefault(extraPayloadRaw, []aiFilterExtraPayloadType{aiFilterExtraPayloadArticleContent})
	placement := parseAIContentProcessPlacement(placementRaw)

	transFunc := func(item *feeds.Item) (string, error) {
		if rule == "" {
			return "", fmt.Errorf("ai-content-process requires rule param")
		}
		original := getPrimaryFeedItemContent(item)
		generated, err := processAIContentForItem(item, rule, payloadTypes)
		if err != nil {
			return "", err
		}
		return applyAIContentProcessPlacement(original, generated, placement), nil
	}

	cachedTransformer := GetCommonCachedTransformer(
		newAIContentProcessCacheKeyGenerator(rule, payloadTypes, placement),
		transFunc,
		"ai-content-process",
	)
	return OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer))
}

func processAIContentForItem(item *feeds.Item, rule string, payloadTypes []aiFilterExtraPayloadType) (string, error) {
	context, err := buildAIContentProcessArticlePayload(item, payloadTypes)
	if err != nil {
		return "", err
	}

	prompt := buildAIContentProcessPrompt(rule)
	result, err := llmContextCaller(prompt, context, util.ContentProcessOption{
		RemoveImage: true,
		ConvertToMd: true,
	})
	if err != nil {
		return "", err
	}
	return normalizeAIContentProcessMarkdown(result), nil
}

func buildAIContentProcessArticlePayload(item *feeds.Item, payloadTypes []aiFilterExtraPayloadType) (string, error) {
	summary := ""
	if loContainsAIFilterPayload(payloadTypes, aiFilterExtraPayloadArticleSummary) {
		generated, err := generateAIFilterArticleSummary(item)
		if err != nil {
			return "", err
		}
		summary = generated
	}
	return buildAIFilterArticlePayload(item, payloadTypes, summary)
}

func buildAIContentProcessPrompt(rule string) string {
	return fmt.Sprintf(`You are an RSS article content processing assistant.

User processing rule:
%s

Process the article according to the rule.

Output requirements:
- Return only the processed article content in Markdown.
- Do not include explanations, greetings, metadata, or unrelated text.
- Do not wrap the answer in markdown code fences.
- Do not use unnecessary code blocks.
- Preserve useful links and images when they are relevant to the processed result.`, rule)
}

func parseAIContentProcessPlacement(raw string) aiContentProcessPlacement {
	switch aiContentProcessPlacement(strings.ToLower(strings.TrimSpace(raw))) {
	case aiContentProcessPlacementReplace:
		return aiContentProcessPlacementReplace
	case aiContentProcessPlacementAppend:
		return aiContentProcessPlacementAppend
	default:
		return aiContentProcessPlacementPrepend
	}
}

func applyAIContentProcessPlacement(originalHTML string, generatedMarkdown string, placement aiContentProcessPlacement) string {
	generatedHTML := util.Markdown2HTML(generatedMarkdown)
	switch placement {
	case aiContentProcessPlacementReplace:
		return generatedHTML
	case aiContentProcessPlacementAppend:
		return fmt.Sprintf(`<div><div>%s</div><hr/><br/><div>%s</div></div>`, originalHTML, generatedHTML)
	default:
		return combineArticleHTMLWithGeneratedMarkdown(originalHTML, generatedMarkdown)
	}
}

func normalizeAIContentProcessMarkdown(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "```") {
		return trimmed
	}

	lines := strings.Split(trimmed, "\n")
	if len(lines) < 3 {
		return trimmed
	}
	infoString := strings.TrimSpace(strings.TrimPrefix(lines[0], "```"))
	if infoString != "" && infoString != "markdown" && infoString != "md" {
		return trimmed
	}
	if strings.TrimSpace(lines[len(lines)-1]) != "```" {
		return trimmed
	}
	return strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))
}

func newAIContentProcessCacheKeyGenerator(rule string, payloadTypes []aiFilterExtraPayloadType, placement aiContentProcessPlacement) ContentCacheKeyGenerator {
	ruleHash := util.GetTextContentHash(rule)
	payloadHash := util.GetTextContentHash(strings.Join(aiFilterPayloadTypesToStrings(payloadTypes), ","))
	return func(item *feeds.Item) (string, error) {
		return util.GetTextContentHash(strings.Join([]string{
			ruleHash,
			payloadHash,
			string(placement),
			strings.TrimSpace(item.Title),
			util.GetTextContentHash(getPrimaryFeedItemContent(item)),
		}, "|")), nil
	}
}

func aiFilterPayloadTypesToStrings(payloadTypes []aiFilterExtraPayloadType) []string {
	values := make([]string, 0, len(payloadTypes))
	for _, payloadType := range payloadTypes {
		values = append(values, string(payloadType))
	}
	return values
}

func loContainsAIFilterPayload(payloadTypes []aiFilterExtraPayloadType, target aiFilterExtraPayloadType) bool {
	for _, payloadType := range payloadTypes {
		if payloadType == target {
			return true
		}
	}
	return false
}
