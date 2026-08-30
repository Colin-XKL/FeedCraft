package craft

import (
	"context"
	"testing"

	"FeedCraft/internal/model"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/require"
)

type nilResultProcessor struct{}

func (nilResultProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	_ = ctx
	_ = feed
	return nil, nil
}

func TestRestoreLegacyItemMetadata_PreservesDistinctDuplicateKeys(t *testing.T) {
	enc1 := &feeds.Enclosure{Url: "https://example.com/a.mp3"}
	enc2 := &feeds.Enclosure{Url: "https://example.com/b.mp3"}
	originals := []*feeds.Item{
		{Id: "same", Link: &feeds.Link{Href: "https://example.com/p"}, Enclosure: enc1, IsPermaLink: "true"},
		{Id: "same", Link: &feeds.Link{Href: "https://example.com/p"}, Enclosure: enc2, IsPermaLink: "false"},
	}
	converted := &feeds.Feed{
		Items: []*feeds.Item{
			{Id: "same", Link: &feeds.Link{Href: "https://example.com/p"}},
			{Id: "same", Link: &feeds.Link{Href: "https://example.com/p"}},
		},
	}

	restoreLegacyItemMetadata(originals, converted)

	require.Equal(t, enc1, converted.Items[0].Enclosure)
	require.Equal(t, "true", converted.Items[0].IsPermaLink)
	require.Equal(t, enc2, converted.Items[1].Enclosure)
	require.Equal(t, "false", converted.Items[1].IsPermaLink)
}

func TestApplyLocalProcessorToLegacyFeed_NilProcessorResultDoesNotPanic(t *testing.T) {
	feed := &feeds.Feed{
		Title: "keep original",
		Items: []*feeds.Item{{Title: "item"}},
	}

	err := applyLocalProcessorToLegacyFeed(context.Background(), nilResultProcessor{}, feed)
	require.Error(t, err)
	require.Contains(t, err.Error(), "nil feed")
	require.Equal(t, "keep original", feed.Title)
	require.Len(t, feed.Items, 1)
}
