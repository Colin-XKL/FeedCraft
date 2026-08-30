package craft

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"FeedCraft/internal/util"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const longBody = "this article body is long enough to exceed the minimum content length required by the llm condition checker"

func TestOptionIgnoreAdvertorial_RemovesArticlesMarkedAsAds(t *testing.T) {
	original := llmContextCaller
	llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
		if strings.Contains(context, "Sponsored Deal") {
			return "true", nil
		}
		return "false", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Real News", Content: longBody},
			{Title: "Sponsored Deal", Content: longBody},
			{Title: "Another Story", Content: longBody},
		},
	}

	option := OptionIgnoreAdvertorial("decide if advertorial")
	require.NoError(t, option(feed, ExtraPayload{}))
	assert.Equal(t, []string{"Real News", "Another Story"}, itemTitles(feed.Items))
}

func TestOptionIgnoreAdvertorial_UsesDescriptionWhenContentEmpty(t *testing.T) {
	var seenContexts []string
	original := llmContextCaller
	llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
		seenContexts = append(seenContexts, context)
		return "false", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Fallback", Description: longBody},
		},
	}

	option := OptionIgnoreAdvertorial("decide if advertorial")
	require.NoError(t, option(feed, ExtraPayload{}))
	require.Len(t, feed.Items, 1)
	require.Len(t, seenContexts, 1)
	assert.Contains(t, seenContexts[0], longBody)
}

func TestOptionIgnoreAdvertorial_KeepsArticleOnLLMError(t *testing.T) {
	original := llmContextCaller
	llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
		return "", fmt.Errorf("llm down")
	}
	t.Cleanup(func() { llmContextCaller = original })

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Survivor", Content: longBody},
		},
	}

	option := OptionIgnoreAdvertorial("decide if advertorial")
	require.NoError(t, option(feed, ExtraPayload{}))
	assert.Equal(t, []string{"Survivor"}, itemTitles(feed.Items))
}

func TestOptionIgnoreAdvertorial_EmptyFeedIsNoop(t *testing.T) {
	feed := &feeds.Feed{}
	option := OptionIgnoreAdvertorial("decide if advertorial")
	require.NoError(t, option(feed, ExtraPayload{}))
	assert.Empty(t, feed.Items)
}

func TestOptionLLMFilterGeneric_RemovesMatchingArticles(t *testing.T) {
	original := llmContextCaller
	llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
		if strings.Contains(context, "Sports") {
			return "true", nil
		}
		return "false", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Sports Update", Content: longBody},
			{Title: "Tech News", Content: longBody},
		},
	}

	option := OptionLLMFilterGeneric("Is this about sports?")
	require.NoError(t, option(feed, ExtraPayload{}))
	assert.Equal(t, []string{"Tech News"}, itemTitles(feed.Items))
}

func TestOptionLLMFilterGeneric_KeepsArticleWhenLLMErrors(t *testing.T) {
	original := llmContextCaller
	llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
		return "", fmt.Errorf("llm error")
	}
	t.Cleanup(func() { llmContextCaller = original })

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Kept On Error", Content: longBody},
		},
	}

	option := OptionLLMFilterGeneric("Is this spam?")
	require.NoError(t, option(feed, ExtraPayload{}))
	assert.Equal(t, []string{"Kept On Error"}, itemTitles(feed.Items))
}
