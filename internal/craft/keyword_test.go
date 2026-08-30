package craft

import (
	"testing"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func keywordTestFeed() *feeds.Feed {
	return &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Golang release notes", Content: "performance improvements", Description: "go runtime"},
			{Title: "Cooking recipes", Content: "delicious pasta", Description: "italian food"},
			{Title: "Weather report", Content: "sunny with golang clouds", Description: "forecast"},
		},
	}
}

func itemTitles(items []*feeds.Item) []string {
	titles := make([]string, 0, len(items))
	for _, item := range items {
		titles = append(titles, item.Title)
	}
	return titles
}

func TestOptionKeyword_IncludeAndExcludeScopes(t *testing.T) {
	tests := []struct {
		name     string
		mode     KeywordFilterMode
		scope    KeywordMatchScope
		keywords []string
		want     []string
	}{
		{
			name:     "include by title only",
			mode:     KeywordIncludeMode,
			scope:    KeywordMatchTitle,
			keywords: []string{"Golang"},
			want:     []string{"Golang release notes"},
		},
		{
			name:     "include by content only matches body",
			mode:     KeywordIncludeMode,
			scope:    KeywordMatchContent,
			keywords: []string{"golang"},
			want:     []string{"Weather report"},
		},
		{
			name:     "include by all matches title or content",
			mode:     KeywordIncludeMode,
			scope:    KeywordMatchAll,
			keywords: []string{"golang", "Golang"},
			want:     []string{"Golang release notes", "Weather report"},
		},
		{
			name:     "exclude by title removes matched",
			mode:     KeywordExcludeMode,
			scope:    KeywordMatchTitle,
			keywords: []string{"Golang"},
			want:     []string{"Cooking recipes", "Weather report"},
		},
		{
			name:     "content scope also matches description",
			mode:     KeywordIncludeMode,
			scope:    KeywordMatchContent,
			keywords: []string{"italian"},
			want:     []string{"Cooking recipes"},
		},
		{
			name:     "no match yields empty for include",
			mode:     KeywordIncludeMode,
			scope:    KeywordMatchAll,
			keywords: []string{"nonexistent"},
			want:     []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			feed := keywordTestFeed()
			option := optionKeyword(tc.mode, tc.scope, tc.keywords)
			require.NoError(t, option(feed, ExtraPayload{}))
			assert.Equal(t, tc.want, itemTitles(feed.Items))
		})
	}
}

func TestOptionKeyword_EmptyKeywordListIsNoop(t *testing.T) {
	feed := keywordTestFeed()
	option := optionKeyword(KeywordIncludeMode, KeywordMatchAll, nil)
	require.NoError(t, option(feed, ExtraPayload{}))
	assert.Len(t, feed.Items, 3)
}

func TestOptionKeyword_UnknownModeDropsEverything(t *testing.T) {
	feed := keywordTestFeed()
	option := optionKeyword(KeywordFilterMode("bogus"), KeywordMatchAll, []string{"Golang"})
	require.NoError(t, option(feed, ExtraPayload{}))
	assert.Empty(t, feed.Items)
}

func TestKeywordCraftLoadParams_ParsesModeAndKeywords(t *testing.T) {
	options := keywordCraftLoadParams(map[string]string{
		"mode":     "exclude",
		"scope":    "title",
		"keywords": "Golang,Cooking",
	})
	require.Len(t, options, 1)

	feed := keywordTestFeed()
	require.NoError(t, options[0](feed, ExtraPayload{}))
	assert.Equal(t, []string{"Weather report"}, itemTitles(feed.Items))
}

// TestKeywordCraftLoadParams_ContentScopeFallsBackToTitle pins the current
// behavior of the legacy loader: scope="content" is mapped to title matching
// (see keyword.go). The native runtime path handles content scope correctly;
// this regression test documents the legacy quirk so a future fix is intentional.
func TestKeywordCraftLoadParams_ContentScopeFallsBackToTitle(t *testing.T) {
	options := keywordCraftLoadParams(map[string]string{
		"mode":     "include",
		"scope":    "content",
		"keywords": "golang",
	})
	require.Len(t, options, 1)

	feed := keywordTestFeed()
	require.NoError(t, options[0](feed, ExtraPayload{}))
	// Lowercase "golang" only appears in content bodies, but because the loader
	// maps content scope to title matching, nothing matches in the title.
	assert.Empty(t, feed.Items)
}
