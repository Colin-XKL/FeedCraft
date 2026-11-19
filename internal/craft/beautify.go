package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"fmt"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

// llmCaller is a variable to allow mocking in tests
var llmCaller = adapter.SimpleLLMCall

const beautifyArticleContentPrompt = `
You are a professional editor. Your task is to reformat the following article content into clean, standard Markdown.
The input content is HTML.
Follow these rules:
1. Preserve the original meaning and wording. Do not summarize or rewrite the content unless necessary for clarity.
2. Fix formatting issues like broken line breaks, excessive whitespace, or messy lists.
3. Remove advertisements, promotional banners, "read more" links, and irrelevant footer info.
4. Ensure images are preserved using standard Markdown '![]()' syntax. Keep the image source URL exactly as is.
5. If there are captions or notes, format them clearly (e.g., using italics or blockquotes).
6. Return ONLY the Markdown content. Do not include any explanations or conversational text.
`

func beautifyCompactedArticleContent(compactedContent string, prompt string) (string, error) {
	if strings.TrimSpace(compactedContent) == "" {
		return "", fmt.Errorf("empty content")
	}

	finalPrompt := fmt.Sprintf("%s\n\n---\n\n%s", prompt, compactedContent)

	beautifiedMd, err := llmCaller(adapter.UseDefaultModel, finalPrompt)
	if err != nil {
		return "", err
	}

	return util.MarkdownToHTML(beautifiedMd), nil
}

func beautifyArticleContent(content string, prompt string) (string, error) {
	// 1. Prepare prompt with original HTML content
	// Check if content is empty
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("empty content")
	}

	compactor := util.NewLLMContentCompactor()
	compactor.Reserve(prompt)
	compactedContent := compactor.Compact(content)

	// 2. Call LLM to beautify the content and convert to Markdown
	renderedHTML, err := beautifyCompactedArticleContent(compactedContent, prompt)
	if err != nil {
		return "", err
	}

	return compactor.Restore(renderedHTML), nil
}

func cachedBeautifyArticleContent(title string, content string, prompt string) (string, error) {
	if strings.TrimSpace(content) == "" {
		return "", nil
	}

	compactor := util.NewLLMContentCompactor()
	compactor.Reserve(prompt)
	compactedContent := compactor.Compact(content)

	hashVal := util.GetTextContentHash(strings.Join([]string{
		util.GetTextContentHash(prompt),
		strings.TrimSpace(title),
		strings.TrimSpace(compactedContent),
	}, "|"))
	cacheKey := getCraftCacheKey("beautify article content", hashVal)

	renderedHTML, err := util.CachedFuncWithPreLog(cacheKey, func() (string, error) {
		return beautifyCompactedArticleContent(compactedContent, prompt)
	}, func(isCached bool) {
		logrus.Infof("applying craft [beautify article content] to article [%s], cached: %v", title, isCached)
	})
	if err != nil {
		return "", err
	}

	return compactor.Restore(renderedHTML), nil
}

// GetBeautifyContentCraftOptions returns the craft options for beautification
func GetBeautifyContentCraftOptions(prompt string) []LegacyCraftOption {
	cachedTransformer := func(item *feeds.Item) (string, error) {
		content := item.Content
		if content == "" {
			content = item.Description
		}
		if strings.TrimSpace(content) == "" {
			return "", fmt.Errorf("empty content")
		}
		return cachedBeautifyArticleContent(item.Title, content, prompt)
	}

	craftOption := []LegacyCraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer)),
	}
	return craftOption
}

func beautifyContentCraftLoadParam(m map[string]string) []LegacyCraftOption {
	prompt, exist := m["prompt"]
	if !exist || len(prompt) == 0 {
		prompt = beautifyArticleContentPrompt
	}
	return GetBeautifyContentCraftOptions(prompt)
}

var beautifyContentParamTmpl = []ParamTemplate{
	{Key: "prompt", Description: "prompt for using llm to beautify content", Default: beautifyArticleContentPrompt},
}
