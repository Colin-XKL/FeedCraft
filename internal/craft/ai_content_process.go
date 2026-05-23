package craft

import (
	"fmt"
	"strings"

	"FeedCraft/internal/util"

	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
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

func aiContentProcessCraftLoadParam(m map[string]string) []LegacyCraftOption {
	return GetAIContentProcessCraftOptions(m["rule"], m["extra-payload"], m["placement"])
}

func GetAIContentProcessCraftOptions(rule string, extraPayloadRaw string, placementRaw string) []LegacyCraftOption {
	return []LegacyCraftOption{
		OptionAIContentProcess(rule, extraPayloadRaw, placementRaw),
	}
}

func OptionAIContentProcess(rule string, extraPayloadRaw string, placementRaw string) LegacyCraftOption {
	rule = strings.TrimSpace(rule)
	payloadTypes := parseAIFilterExtraPayloadWithDefault(extraPayloadRaw, []aiFilterExtraPayloadType{aiFilterExtraPayloadArticleContent})
	placement := parseAIContentProcessPlacement(placementRaw)
	prompt := buildAIContentProcessPrompt(rule)

	transFunc := func(item *feeds.Item) (string, error) {
		if rule == "" {
			return "", fmt.Errorf("ai-content-process requires rule param")
		}
		return cachedAIContentProcessItem(item, prompt, payloadTypes, placement)
	}

	return OptionTransformFeedItem(GetArticleContentProcessor(transFunc))
}

func cachedAIContentProcessItem(item *feeds.Item, prompt string, payloadTypes []aiFilterExtraPayloadType, placement aiContentProcessPlacement) (string, error) {
	original := getPrimaryFeedItemContent(item)
	context, err := buildAIContentProcessArticlePayload(item, payloadTypes)
	if err != nil {
		return "", err
	}
	cacheKey := getCraftCacheKey("ai-content-process-result", util.GetTextContentHash(strings.Join([]string{
		util.GetTextContentHash(prompt),
		util.GetTextContentHash(context),
		string(placement),
		util.GetTextContentHash(original),
	}, "|")))

	return util.CachedFuncWithPreLog(cacheKey, func() (string, error) {
		result, callErr := llmContextCaller(prompt, context, util.ContentProcessOption{
			RemoveImage: true,
			ConvertToMd: true,
			Temperature: util.LowestLLMTemperaturePtr(),
		})
		if callErr != nil {
			return "", callErr
		}
		generated := normalizeAIContentProcessMarkdown(result)
		return applyAIContentProcessPlacement(original, generated, placement), nil
	}, func(isCached bool) {
		title := ""
		if item != nil {
			title = item.Title
		}
		logrus.Infof("applying craft [ai-content-process] to article [%s], cached: %v", title, isCached)
	})
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
		return combineArticleHTMLFragments(originalHTML, generatedHTML)
	default:
		return combineArticleHTMLFragments(generatedHTML, originalHTML)
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
	infoString := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(lines[0], "```")))
	if infoString != "" && infoString != "markdown" && infoString != "md" {
		return trimmed
	}
	if strings.TrimSpace(lines[len(lines)-1]) != "```" {
		return trimmed
	}
	return strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))
}

func loContainsAIFilterPayload(payloadTypes []aiFilterExtraPayloadType, target aiFilterExtraPayloadType) bool {
	for _, payloadType := range payloadTypes {
		if payloadType == target {
			return true
		}
	}
	return false
}
