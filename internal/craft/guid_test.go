package craft

import (
	"testing"

	"FeedCraft/internal/util"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedItemGuidGenerator_DeterministicForSameContent(t *testing.T) {
	itemA := &feeds.Item{Title: "Hello", Content: "World", Description: "desc"}
	itemB := &feeds.Item{Title: "Hello", Content: "World", Description: "desc"}

	guidA, err := feedItemGuidGenerator(itemA)
	require.NoError(t, err)
	guidB, err := feedItemGuidGenerator(itemB)
	require.NoError(t, err)

	assert.Equal(t, guidA, guidB)
	assert.Equal(t, util.GetTextContentHash("HelloWorlddesc"), guidA)
}

func TestFeedItemGuidGenerator_DiffersWhenContentDiffers(t *testing.T) {
	guidA, err := feedItemGuidGenerator(&feeds.Item{Title: "Hello", Content: "World"})
	require.NoError(t, err)
	guidB, err := feedItemGuidGenerator(&feeds.Item{Title: "Hello", Content: "Mars"})
	require.NoError(t, err)

	assert.NotEqual(t, guidA, guidB)
}

func TestFeedItemGuidGenerator_EmptyFieldsYieldUniqueRandomGuids(t *testing.T) {
	guidA, err := feedItemGuidGenerator(&feeds.Item{})
	require.NoError(t, err)
	guidB, err := feedItemGuidGenerator(&feeds.Item{})
	require.NoError(t, err)

	assert.NotEmpty(t, guidA)
	assert.NotEmpty(t, guidB)
	// With no title/content/description, a random uuid is returned each call.
	assert.NotEqual(t, guidA, guidB)
}

func TestGetGuidProcessor_OverwritesItemId(t *testing.T) {
	processor := GetGuidProcessor(func(item *feeds.Item) (string, error) {
		return "computed-id", nil
	})

	item := &feeds.Item{Id: "stale", Title: "title"}
	require.NoError(t, processor(item, ExtraPayload{}))
	assert.Equal(t, "computed-id", item.Id)
}

func TestGetGuidCraftOptions_RewritesIdsAcrossFeed(t *testing.T) {
	options := GetGuidCraftOptions()
	require.Len(t, options, 1)

	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Id: "old-1", Title: "First", Content: "alpha"},
			{Id: "old-2", Title: "Second", Content: "beta"},
		},
	}
	require.NoError(t, options[0](feed, ExtraPayload{}))

	assert.Equal(t, util.GetTextContentHash("Firstalpha"), feed.Items[0].Id)
	assert.Equal(t, util.GetTextContentHash("Secondbeta"), feed.Items[1].Id)
	assert.NotEqual(t, feed.Items[0].Id, feed.Items[1].Id)
}
