package craft

import (
	"FeedCraft/internal/adapter"
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/feeds"
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

func beautifyArticleContent(content string, prompt string) (string, error) {
	// 1. Prepare prompt with original HTML content
	// Check if content is empty
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("empty content")
	}

	finalPrompt := fmt.Sprintf("%s\n\n---\n\n%s", prompt, content)

	// 2. Call LLM to beautify the content and convert to Markdown
	beautifiedMd, err := llmCaller(adapter.UseDefaultModel, finalPrompt)
	if err != nil {
		return "", err
	}

	// 3. Convert Beautified Markdown back to HTML
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(beautifiedMd))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	beautifiedHtml := markdown.Render(doc, renderer)

	return string(beautifiedHtml), nil
}

// GetBeautifyContentCraftOptions returns the craft options for beautification
func GetBeautifyContentCraftOptions(prompt string) []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		content := item.Content
		if content == "" {
			content = item.Description
		}
		return beautifyArticleContent(content, prompt)
	}

	cachedTransformer := GetCommonCachedTransformer(
		cacheKeyForArticleContent, transFunc, "beautify article content")

	craftOption := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer)),
	}
	return craftOption
}

func beautifyContentCraftLoadParam(m map[string]string) []CraftOption {
	prompt, exist := m["prompt"]
	if !exist || len(prompt) == 0 {
		prompt = beautifyArticleContentPrompt
	}
	return GetBeautifyContentCraftOptions(prompt)
}

var beautifyContentParamTmpl = []ParamTemplate{
	{Key: "prompt", Description: "prompt for using llm to beautify content", Default: beautifyArticleContentPrompt},
}
