package craft

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"FeedCraft/internal/util"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAIFilterDecisionAcceptsFencedJSON(t *testing.T) {
	decision, err := parseAIFilterDecision("```json\n{\"reason\":\"not relevant\",\"result\":\"drop\"}\n```")

	require.NoError(t, err)
	assert.Equal(t, aiFilterResultDrop, decision.Result)
	assert.Equal(t, "not relevant", decision.Reason)
}

func TestBuildAIFilterArticlePayloadIncludesRequestedContent(t *testing.T) {
	item := &feeds.Item{
		Title:       "AI chip news",
		Description: "short description",
		Content:     "<p>complete article content</p>",
	}

	payload, err := buildAIFilterArticlePayload(item, []aiFilterExtraPayloadType{aiFilterExtraPayloadArticleContent}, "")

	require.NoError(t, err)
	assert.Contains(t, payload, "Article Title:")
	assert.Contains(t, payload, "AI chip news")
	assert.Contains(t, payload, "Article Content:")
	assert.Contains(t, payload, "complete article content")
	assert.NotContains(t, payload, "short description")
}

func TestOptionAIFilterDropsOnlyDropDecisionAndUsesSummaryPayload(t *testing.T) {
	setupTestRedis(t)

	original := llmContextCaller
	var filterContexts []string
	var summaryContexts []string
	var seenMu sync.Mutex
	llmContextCaller = func(prompt, context string, option util.ContentProcessOption) (string, error) {
		if strings.Contains(prompt, "professional summarizer") {
			seenMu.Lock()
			summaryContexts = append(summaryContexts, context)
			seenMu.Unlock()
			if strings.Contains(context, "drop article original") {
				return "summary says drop", nil
			}
			return "summary says keep", nil
		}

		seenMu.Lock()
		filterContexts = append(filterContexts, context)
		seenMu.Unlock()
		if strings.Contains(context, "summary says drop") {
			return `{"reason":"summary matched exclusion","result":"drop"}`, nil
		}
		return `{"reason":"summary did not match exclusion","result":"keep"}`, nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Drop", Content: "<p>drop article original</p>"},
			{Title: "Keep", Content: "<p>keep article original</p>"},
		},
	}

	err := OptionAIFilter("只保留科技有关的文章", "article_summary")(feed, ExtraPayload{})

	require.NoError(t, err)
	require.Len(t, feed.Items, 1)
	assert.Equal(t, "Keep", feed.Items[0].Title)
	require.Len(t, summaryContexts, 2)
	require.Len(t, filterContexts, 2)
	combinedFilterContexts := strings.Join(filterContexts, "\n---\n")
	assert.Contains(t, combinedFilterContexts, "Article Summary:")
	assert.Contains(t, combinedFilterContexts, "summary says drop")
	assert.NotContains(t, combinedFilterContexts, "drop article original")
}

func TestOptionAIFilterKeepsArticleOnInvalidLLMResponse(t *testing.T) {
	setupTestRedis(t)

	original := llmContextCaller
	llmContextCaller = func(prompt, context string, option util.ContentProcessOption) (string, error) {
		return "not json", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Keep on malformed response", Content: "<p>article content long enough</p>"},
		},
	}

	err := OptionAIFilter("只保留科技有关的文章", "article_content")(feed, ExtraPayload{})

	require.NoError(t, err)
	require.Len(t, feed.Items, 1)
	assert.Equal(t, "Keep on malformed response", feed.Items[0].Title)
}

func TestAIFilterCraftLoadParamUsesRuleParam(t *testing.T) {
	setupTestRedis(t)

	original := llmContextCaller
	llmContextCaller = func(prompt, context string, option util.ContentProcessOption) (string, error) {
		assert.Contains(t, prompt, "只保留科技有关的文章")
		return `{"reason":"not a tech article","result":"drop"}`, nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	options := aiFilterCraftLoadParam(map[string]string{
		"rule":          "只保留科技有关的文章",
		"extra-payload": "article_content",
	})
	require.Len(t, options, 1)

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Sports", Content: "<p>football news</p>"},
		},
	}

	err := options[0](feed, ExtraPayload{})

	require.NoError(t, err)
	require.Empty(t, feed.Items)
}

func TestEvaluateAIFilterItemCachesDecision(t *testing.T) {
	setupTestRedis(t)

	original := llmContextCaller
	filterCalls := 0
	llmContextCaller = func(prompt, context string, option util.ContentProcessOption) (string, error) {
		filterCalls += 1
		return `{"reason":"cached decision","result":"keep"}`, nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	item := &feeds.Item{
		Title:   "Cache Me",
		Content: "<p>same article content</p>",
	}
	payloadTypes := []aiFilterExtraPayloadType{aiFilterExtraPayloadArticleContent}

	first, err := evaluateAIFilterItem(item, "只保留科技有关的文章", payloadTypes)
	require.NoError(t, err)
	second, err := evaluateAIFilterItem(item, "只保留科技有关的文章", payloadTypes)
	require.NoError(t, err)

	assert.Equal(t, aiFilterResultKeep, first.Result)
	assert.Equal(t, first, second)
	assert.Equal(t, 1, filterCalls)
}

func TestBuildAIFilterArticlePayloadIncludesRawRSSItemAsJSON(t *testing.T) {
	item := &feeds.Item{
		Title:       "Raw item title",
		Description: "Raw item description",
		Content:     "<p>Raw item content</p>",
		Id:          "guid-1",
		Link:        &feeds.Link{Href: "https://example.com/post"},
		Created:     time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC),
	}

	payload, err := buildAIFilterArticlePayload(item, []aiFilterExtraPayloadType{aiFilterExtraPayloadRawRSSItem}, "")

	require.NoError(t, err)
	assert.Contains(t, payload, "Raw RSS Item JSON:")
	assert.Contains(t, payload, `"title":"Raw item title"`)
	assert.Contains(t, payload, `"link":"https://example.com/post"`)
	assert.Contains(t, payload, fmt.Sprintf(`"created":%q`, item.Created.Format(time.RFC3339)))
}
