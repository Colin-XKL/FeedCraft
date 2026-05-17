package craft

import (
	"strings"
	"testing"

	"FeedCraft/internal/util"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionAIContentProcessPrependsGeneratedMarkdownAsHTML(t *testing.T) {
	setupTestRedis(t)

	original := llmContextCaller
	var seenPrompt string
	var seenContext string
	llmContextCaller = func(prompt, context string, option util.ContentProcessOption) (string, error) {
		seenPrompt = prompt
		seenContext = context
		assert.True(t, option.ConvertToMd)
		assert.True(t, option.RemoveImage)
		return "## AI Takeaway\n\n- Generated insight", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{
				Title:   "Processor input " + t.Name(),
				Content: "<p>original article content</p>",
			},
		},
	}

	err := OptionAIContentProcess("提取适合忙人阅读的要点", "article_content", string(aiContentProcessPlacementPrepend))(feed, ExtraPayload{})

	require.NoError(t, err)
	require.Len(t, feed.Items, 1)
	content := feed.Items[0].Content
	assert.Contains(t, content, "<h2")
	assert.Contains(t, content, "AI Takeaway")
	assert.Contains(t, content, "<li>Generated insight</li>")
	assert.Contains(t, content, "<p>original article content</p>")
	assert.Less(t, strings.Index(content, "AI Takeaway"), strings.Index(content, "original article content"))
	assert.Equal(t, feed.Items[0].Content, feed.Items[0].Description)
	assert.Contains(t, seenPrompt, "提取适合忙人阅读的要点")
	assert.Contains(t, seenPrompt, "Return only the processed article content in Markdown")
	assert.Contains(t, seenPrompt, "Do not wrap the answer in markdown code fences")
	assert.Contains(t, seenContext, "Article Content:")
	assert.Contains(t, seenContext, "original article content")
}

func TestOptionAIContentProcessSupportsReplaceAndAppendPlacements(t *testing.T) {
	setupTestRedis(t)

	original := llmContextCaller
	llmContextCaller = func(prompt, context string, option util.ContentProcessOption) (string, error) {
		return "### Generated\n\nreplacement or appendix", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	replaceFeed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Replace " + t.Name(), Content: "<p>original body</p>"},
		},
	}
	err := OptionAIContentProcess("改写正文", "article_content", string(aiContentProcessPlacementReplace))(replaceFeed, ExtraPayload{})
	require.NoError(t, err)
	assert.Contains(t, replaceFeed.Items[0].Content, "<h3")
	assert.Contains(t, replaceFeed.Items[0].Content, "replacement or appendix")
	assert.NotContains(t, replaceFeed.Items[0].Content, "original body")

	appendFeed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Append " + t.Name(), Content: "<p>original body</p>"},
		},
	}
	err = OptionAIContentProcess("生成相关阅读建议", "article_content", string(aiContentProcessPlacementAppend))(appendFeed, ExtraPayload{})
	require.NoError(t, err)
	content := appendFeed.Items[0].Content
	assert.Contains(t, content, "<p>original body</p>")
	assert.Contains(t, content, "replacement or appendix")
	assert.Less(t, strings.Index(content, "original body"), strings.Index(content, "replacement or appendix"))
}

func TestOptionAIContentProcessCacheKeyIncludesPlacement(t *testing.T) {
	setupTestRedis(t)

	original := llmContextCaller
	llmContextCaller = func(prompt, context string, option util.ContentProcessOption) (string, error) {
		return "Generated note", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	replaceFeed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Same cache key " + t.Name(), Content: "<p>same original body</p>"},
		},
	}
	err := OptionAIContentProcess("同一规则", "article_content", string(aiContentProcessPlacementReplace))(replaceFeed, ExtraPayload{})
	require.NoError(t, err)
	assert.NotContains(t, replaceFeed.Items[0].Content, "same original body")

	appendFeed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Same cache key " + t.Name(), Content: "<p>same original body</p>"},
		},
	}
	err = OptionAIContentProcess("同一规则", "article_content", string(aiContentProcessPlacementAppend))(appendFeed, ExtraPayload{})
	require.NoError(t, err)
	assert.Contains(t, appendFeed.Items[0].Content, "same original body")
	assert.Less(t, strings.Index(appendFeed.Items[0].Content, "same original body"), strings.Index(appendFeed.Items[0].Content, "Generated note"))
}

func TestAIContentProcessCraftLoadParamUsesDefaultsAndRegistersTemplate(t *testing.T) {
	setupTestRedis(t)

	original := llmContextCaller
	var seenContext string
	llmContextCaller = func(prompt, context string, option util.ContentProcessOption) (string, error) {
		seenContext = context
		return "Generated default summary", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	tmpl, ok := GetSysCraftTemplateDict()["ai-content-process"]
	require.True(t, ok)
	assert.Equal(t, "ai-content-process", tmpl.Name)
	assert.NotEmpty(t, tmpl.ParamTemplateDefine)

	options := aiContentProcessCraftLoadParam(map[string]string{
		"rule": "生成一段摘要",
	})
	require.Len(t, options, 1)

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Default " + t.Name(), Content: "<p>default content payload</p>"},
		},
	}
	err := options[0](feed, ExtraPayload{})
	require.NoError(t, err)
	assert.Contains(t, feed.Items[0].Content, "Generated default summary")
	assert.Contains(t, feed.Items[0].Content, "default content payload")
	assert.Contains(t, seenContext, "Article Content:")
}

func TestNormalizeAIContentProcessMarkdownRemovesCodeFence(t *testing.T) {
	result := normalizeAIContentProcessMarkdown("```markdown\n# Title\n\nBody\n```")

	assert.Equal(t, "# Title\n\nBody", result)
}

func TestNormalizeAIContentProcessMarkdownKeepsIntentionalCodeBlock(t *testing.T) {
	result := normalizeAIContentProcessMarkdown("```go\nfmt.Println(\"hello\")\n```")

	assert.Equal(t, "```go\nfmt.Println(\"hello\")\n```", result)
}
