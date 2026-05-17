package craft

import (
	"strings"
	"testing"
	"time"

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
		require.NotNil(t, option.Temperature)
		assert.Equal(t, 0.0, *option.Temperature)
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

func TestApplyAIContentProcessPlacementOmitsSeparatorWhenOriginalEmpty(t *testing.T) {
	generatedMarkdown := "## Generated\n\nNo original content"

	prepended := applyAIContentProcessPlacement("", generatedMarkdown, aiContentProcessPlacementPrepend)
	appended := applyAIContentProcessPlacement("", generatedMarkdown, aiContentProcessPlacementAppend)

	for _, content := range []string{prepended, appended} {
		assert.Contains(t, content, "<h2")
		assert.Contains(t, content, "No original content")
		assert.NotContains(t, content, "<hr")
		assert.NotContains(t, content, "<br")
	}
	assert.Equal(t, prepended, appended)
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

func TestOptionAIContentProcessCacheKeyIncludesArticleDatePayload(t *testing.T) {
	setupTestRedis(t)

	original := llmContextCaller
	calls := 0
	llmContextCaller = func(prompt, context string, option util.ContentProcessOption) (string, error) {
		calls += 1
		if strings.Contains(context, "2026-05-02") {
			return "Generated for May 2", nil
		}
		return "Generated for May 1", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	firstFeed := &feeds.Feed{
		Items: []*feeds.Item{
			{
				Title:   "Same date payload " + t.Name(),
				Content: "<p>same body</p>",
				Created: time.Date(2026, 5, 1, 10, 0, 0, 0,
					time.UTC),
			},
		},
	}
	err := OptionAIContentProcess("按日期生成说明", "article_date", string(aiContentProcessPlacementReplace))(firstFeed, ExtraPayload{})
	require.NoError(t, err)
	assert.Contains(t, firstFeed.Items[0].Content, "Generated for May 1")

	secondFeed := &feeds.Feed{
		Items: []*feeds.Item{
			{
				Title:   "Same date payload " + t.Name(),
				Content: "<p>same body</p>",
				Created: time.Date(2026, 5, 2, 10, 0, 0, 0,
					time.UTC),
			},
		},
	}
	err = OptionAIContentProcess("按日期生成说明", "article_date", string(aiContentProcessPlacementReplace))(secondFeed, ExtraPayload{})
	require.NoError(t, err)
	assert.Contains(t, secondFeed.Items[0].Content, "Generated for May 2")
	assert.Equal(t, 2, calls)
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

func TestParameterizedSystemCraftTemplatesAreMarkedTemplateOnly(t *testing.T) {
	templates := GetSysCraftTemplateDict()

	assert.True(t, templates["ai-filter"].TemplateOnly)
	assert.True(t, templates["ai-content-process"].TemplateOnly)
	assert.False(t, templates["summary"].TemplateOnly)
}

func TestNormalizeAIContentProcessMarkdownRemovesCodeFence(t *testing.T) {
	result := normalizeAIContentProcessMarkdown("```markdown\n# Title\n\nBody\n```")

	assert.Equal(t, "# Title\n\nBody", result)
}

func TestNormalizeAIContentProcessMarkdownRemovesCaseInsensitiveMarkdownFence(t *testing.T) {
	result := normalizeAIContentProcessMarkdown("```Markdown\n# Title\n\nBody\n```")

	assert.Equal(t, "# Title\n\nBody", result)
}

func TestNormalizeAIContentProcessMarkdownKeepsIntentionalCodeBlock(t *testing.T) {
	result := normalizeAIContentProcessMarkdown("```go\nfmt.Println(\"hello\")\n```")

	assert.Equal(t, "```go\nfmt.Println(\"hello\")\n```", result)
}
