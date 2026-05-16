package craft

import (
	"testing"
	"time"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionLimit_SortsByCreatedTimeBeforeTruncating(t *testing.T) {
	now := time.Now()
	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Id: "oldest", Created: now.Add(-3 * time.Hour)},
			{Id: "newest", Created: now},
			{Id: "middle", Created: now.Add(-1 * time.Hour)},
		},
	}

	option := OptionLimit(2)
	err := option(feed, ExtraPayload{})

	require.NoError(t, err)
	require.Len(t, feed.Items, 2)
	assert.Equal(t, "newest", feed.Items[0].Id)
	assert.Equal(t, "middle", feed.Items[1].Id)
}

func TestOptionLimit_UsesUpdatedTimeWhenCreatedTimeMatches(t *testing.T) {
	now := time.Now()
	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Id: "older-update", Created: now, Updated: now.Add(-3 * time.Hour)},
			{Id: "newer-update", Created: now, Updated: now.Add(-1 * time.Hour)},
			{Id: "oldest-update", Created: now, Updated: now.Add(-6 * time.Hour)},
		},
	}

	option := OptionLimit(2)
	err := option(feed, ExtraPayload{})

	require.NoError(t, err)
	require.Len(t, feed.Items, 2)
	assert.Equal(t, "newer-update", feed.Items[0].Id)
	assert.Equal(t, "older-update", feed.Items[1].Id)
}
