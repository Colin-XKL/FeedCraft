package craft

import (
	"fmt"
	"strings"

	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"

	"github.com/gorilla/feeds"
)

const retitleNoChangeSentinel = "__FEEDCRAFT_KEEP_ORIGINAL_TITLE__"

func renderRetitlePrompt(prompt string) string {
	finalPrompt := renderTargetLangPrompt(prompt, constant.DefaultPrompts[constant.ProcessorTypeRetitle])
	if strings.TrimSpace(prompt) == "" {
		return finalPrompt
	}
	return fmt.Sprintf(`%s

Important output contract:
If the article content is empty or does not contain enough information to create a reliable new title, return exactly: %s
Return only the new title or this special value. Do not include quotes, markdown, labels, or explanations.`, finalPrompt, retitleNoChangeSentinel)
}

func shouldKeepOriginalTitle(generated string) bool {
	trimmed := strings.TrimSpace(generated)
	trimmed = strings.Trim(trimmed, "`\"' \n\t")
	return trimmed == "" || strings.EqualFold(trimmed, retitleNoChangeSentinel)
}

func getFeedsItemPrimaryContent(item *feeds.Item) string {
	if item == nil {
		return ""
	}
	content := item.Content
	if strings.TrimSpace(content) == "" {
		content = item.Description
	}
	return content
}

func getFeedsItemContentForPrompt(item *feeds.Item, original string) string {
	if item == nil {
		return original
	}
	link := ""
	if item.Link != nil {
		link = item.Link.Href
	}
	domain, _ := util.ParseDomainFromUrl(link)
	cleaned := util.HTMLToMarkdown(original, domain)
	if strings.TrimSpace(cleaned) != "" {
		return cleaned
	}
	return original
}

func GetRetitleCraftOptions(prompt string) []LegacyCraftOption {
	finalPrompt := renderRetitlePrompt(prompt)
	promptHash := util.GetTextContentHash(finalPrompt)
	cacheKeyGenerator := func(item *feeds.Item) (string, error) {
		if item == nil {
			return promptHash, nil
		}
		payloadHash := util.GetTextContentHash(strings.Join([]string{
			promptHash,
			strings.TrimSpace(item.Title),
			strings.TrimSpace(getFeedsItemPrimaryContent(item)),
		}, "|"))
		return payloadHash, nil
	}
	transFunc := func(item *feeds.Item) (string, error) {
		if item == nil {
			return "", nil
		}
		originalTitle := item.Title
		originalContent := getFeedsItemPrimaryContent(item)
		if strings.TrimSpace(originalContent) == "" {
			return originalTitle, nil
		}
		contentForPrompt := getFeedsItemContentForPrompt(item, originalContent)
		generated, err := CallLLMForArticleTransform(finalPrompt, originalTitle, contentForPrompt, util.ContentProcessOption{})
		if err != nil {
			return "", err
		}
		if shouldKeepOriginalTitle(generated) {
			return originalTitle, nil
		}
		return strings.TrimSpace(generated), nil
	}
	transformer := GetCommonCachedTransformer(cacheKeyGenerator, transFunc, string(constant.ProcessorTypeRetitle))
	return []LegacyCraftOption{
		OptionTransformFeedItem(GetArticleTitleProcessor(transformer)),
	}
}

func retitleCraftLoadParam(m map[string]string) []LegacyCraftOption {
	prompt := m["prompt"]
	return GetRetitleCraftOptions(prompt)
}

var retitleCraftParamTmpl = []ParamTemplate{
	{
		Key:         "prompt",
		Description: "the llm prompt for regenerating article title",
		Default:     constant.DefaultPrompts[constant.ProcessorTypeRetitle],
	},
}
