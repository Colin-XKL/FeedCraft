package craft

import (
	"testing"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
)

// --- 参数解析测试 ---

func TestEmbeddingFilterLoadParam_EmptyAnchors(t *testing.T) {
	params := map[string]string{
		"anchors": "",
	}
	options := embeddingFilterLoadParam(params)
	assert.Empty(t, options, "empty anchors should return no options")
}

func TestEmbeddingFilterLoadParam_ValidAnchors(t *testing.T) {
	params := map[string]string{
		"anchors":   "人工智能\n机器学习\n深度学习",
		"threshold": "0.7",
	}
	options := embeddingFilterLoadParam(params)
	assert.Len(t, options, 1, "should return exactly one CraftOption")
}

func TestEmbeddingFilterLoadParam_AnchorsWithBlankLines(t *testing.T) {
	params := map[string]string{
		"anchors": "人工智能\n\n  \n机器学习\n",
	}
	options := embeddingFilterLoadParam(params)
	assert.Len(t, options, 1, "blank lines should be filtered out, valid anchors remain")
}

func TestEmbeddingFilterLoadParam_InvalidThreshold(t *testing.T) {
	params := map[string]string{
		"anchors":   "test anchor",
		"threshold": "invalid",
	}
	options := embeddingFilterLoadParam(params)
	assert.Len(t, options, 1, "invalid threshold should fallback to default, still return option")
}

func TestEmbeddingFilterLoadParam_ThresholdOutOfRange(t *testing.T) {
	params := map[string]string{
		"anchors":   "test anchor",
		"threshold": "1.5",
	}
	options := embeddingFilterLoadParam(params)
	assert.Len(t, options, 1, "out-of-range threshold should fallback to default")
}

func TestEmbeddingFilterLoadParam_NegativeThreshold(t *testing.T) {
	params := map[string]string{
		"anchors":   "test anchor",
		"threshold": "-0.1",
	}
	options := embeddingFilterLoadParam(params)
	assert.Len(t, options, 1, "negative threshold should fallback to default")
}

func TestEmbeddingFilterLoadParam_DefaultValues(t *testing.T) {
	params := map[string]string{
		"anchors": "test anchor",
	}
	options := embeddingFilterLoadParam(params)
	assert.Len(t, options, 1, "should use default threshold and max_content_length")
}

func TestEmbeddingFilterLoadParam_InvalidMaxContentLength(t *testing.T) {
	params := map[string]string{
		"anchors":            "test anchor",
		"max_content_length": "abc",
	}
	options := embeddingFilterLoadParam(params)
	assert.Len(t, options, 1, "invalid max_content_length should fallback to default")
}

func TestEmbeddingFilterLoadParam_ZeroMaxContentLength(t *testing.T) {
	params := map[string]string{
		"anchors":            "test anchor",
		"max_content_length": "0",
	}
	options := embeddingFilterLoadParam(params)
	assert.Len(t, options, 1, "zero max_content_length should fallback to default")
}

// --- buildArticleText 测试 ---

func TestBuildArticleText_TitleAndContent(t *testing.T) {
	item := &feeds.Item{
		Title:   "Test Title",
		Content: "Test Content Body",
	}
	result := buildArticleText(item, 2000)
	assert.Equal(t, "Test Title\nTest Content Body", result)
}

func TestBuildArticleText_TitleOnly(t *testing.T) {
	item := &feeds.Item{
		Title: "Test Title",
	}
	result := buildArticleText(item, 2000)
	assert.Equal(t, "Test Title", result)
}

func TestBuildArticleText_ContentOnly(t *testing.T) {
	item := &feeds.Item{
		Content: "Test Content Body",
	}
	result := buildArticleText(item, 2000)
	assert.Equal(t, "Test Content Body", result)
}

func TestBuildArticleText_FallbackToDescription(t *testing.T) {
	item := &feeds.Item{
		Title:       "Test Title",
		Description: "Test Description",
	}
	result := buildArticleText(item, 2000)
	assert.Equal(t, "Test Title\nTest Description", result)
}

func TestBuildArticleText_ContentTruncation(t *testing.T) {
	longContent := ""
	for i := 0; i < 300; i++ {
		longContent += "abcdefghij" // 3000 chars total
	}
	item := &feeds.Item{
		Title:   "Title",
		Content: longContent,
	}
	result := buildArticleText(item, 100)
	// 标题 + "\n" + 截取后的100字符
	assert.Equal(t, "Title\n"+longContent[:100], result)
}

func TestBuildArticleText_EmptyItem(t *testing.T) {
	item := &feeds.Item{}
	result := buildArticleText(item, 2000)
	assert.Equal(t, "", result)
}

// --- OptionEmbeddingFilter 边界测试 ---

func TestOptionEmbeddingFilter_EmptyItems(t *testing.T) {
	option := OptionEmbeddingFilter([]string{"test"}, 0.6, 2000, "")
	feed := &feeds.Feed{Items: []*feeds.Item{}}
	err := option(feed, ExtraPayload{})
	assert.NoError(t, err)
	assert.Empty(t, feed.Items)
}

func TestOptionEmbeddingFilter_EmptyAnchors(t *testing.T) {
	option := OptionEmbeddingFilter([]string{}, 0.6, 2000, "")
	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Article 1", Content: "Content 1"},
			{Title: "Article 2", Content: "Content 2"},
		},
	}
	err := option(feed, ExtraPayload{})
	assert.NoError(t, err)
	// 锚点为空时不过滤，保留所有文章
	assert.Len(t, feed.Items, 2)
}

// --- CraftTemplate 注册测试 ---

func TestEmbeddingFilterTemplateRegistered(t *testing.T) {
	templates := GetSysCraftTemplateDict()
	tmpl, exists := templates["embedding-filter"]
	assert.True(t, exists, "embedding-filter template should be registered")
	assert.Equal(t, "embedding-filter", tmpl.Name)
	assert.NotEmpty(t, tmpl.Description)
	assert.NotEmpty(t, tmpl.ParamTemplateDefine)

	// 验证参数模板包含所有必要字段
	paramKeys := make(map[string]bool)
	for _, p := range tmpl.ParamTemplateDefine {
		paramKeys[p.Key] = true
	}
	assert.True(t, paramKeys["anchors"], "should have 'anchors' param")
	assert.True(t, paramKeys["threshold"], "should have 'threshold' param")
	assert.True(t, paramKeys["max_content_length"], "should have 'max_content_length' param")
	assert.True(t, paramKeys["instruction"], "should have 'instruction' param")
}
