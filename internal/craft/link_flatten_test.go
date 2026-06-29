package craft

import (
	"context"
	"testing"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/model"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkFlattenDefaultsKeepOriginalAndExtractExternalLinksOnly(t *testing.T) {
	db := newCraftRuntimeTestDB(t)

	processor, err := BuildOptionChain(db, "link-flatten", "https://example.com/feed.xml")
	require.NoError(t, err)

	feed := &model.CraftFeed{
		Title: "Feed",
		Link:  "https://example.com",
		Articles: []*model.CraftArticle{
			{
				Title:   "Source Article",
				Link:    "https://example.com/posts/1",
				Content: `<p>Read <a href="https://external.test/a">External A</a>, <a href="/internal">Internal</a>, and <a href="mailto:hello@example.com">Mail</a>.</p>`,
			},
		},
	}

	result, err := processor(context.Background(), feed)
	require.NoError(t, err)

	require.Len(t, result.Articles, 2)
	assert.Equal(t, "Source Article", result.Articles[0].Title)
	assert.Equal(t, "https://example.com/posts/1", result.Articles[0].Link)
	assert.Equal(t, "External A", result.Articles[1].Title)
	assert.Equal(t, "https://external.test/a", result.Articles[1].Link)
	assert.Contains(t, result.Articles[1].Description, "Source Article")
	assert.Contains(t, result.Articles[1].Description, "https://example.com/posts/1")
}

func TestLinkFlattenCanDropOriginalAndIncludeSameDomainLinks(t *testing.T) {
	db := newCraftRuntimeTestDB(t)
	require.NoError(t, dao.CreateCraftAtom(db, &dao.CraftAtom{
		Name:         "link-flatten-all",
		TemplateName: "link-flatten",
		Params: map[string]string{
			"keep-original": "false",
			"external-only": "false",
		},
	}))

	processor, err := BuildOptionChain(db, "link-flatten-all", "https://example.com/feed.xml")
	require.NoError(t, err)

	feed := &model.CraftFeed{
		Title: "Feed",
		Link:  "https://example.com",
		Articles: []*model.CraftArticle{
			{
				Title:       "Source Article",
				Link:        "https://example.com/posts/1",
				Description: `<a href="/same">Same Domain</a><a href="https://external.test/b">External B</a>`,
			},
		},
	}

	result, err := processor(context.Background(), feed)
	require.NoError(t, err)

	require.Len(t, result.Articles, 2)
	assert.Equal(t, "Same Domain", result.Articles[0].Title)
	assert.Equal(t, "https://example.com/same", result.Articles[0].Link)
	assert.Equal(t, "External B", result.Articles[1].Title)
	assert.Equal(t, "https://external.test/b", result.Articles[1].Link)
}

func TestLinkFlattenLegacyOptionUsesTemplateDefaults(t *testing.T) {
	options := GetSysCraftTemplateDict()["link-flatten"].GetOptions(map[string]string{})
	require.Len(t, options, 1)

	feed := &feeds.Feed{
		Link: &feeds.Link{Href: "https://example.com"},
		Items: []*feeds.Item{
			{
				Title:       "Legacy Source",
				Link:        &feeds.Link{Href: "https://example.com/post"},
				Description: `<a href="https://elsewhere.test/path">Elsewhere</a><a href="/same">Same</a>`,
			},
		},
	}

	err := options[0](feed, ExtraPayload{originalFeedUrl: "https://example.com/feed.xml"})
	require.NoError(t, err)

	require.Len(t, feed.Items, 2)
	assert.Equal(t, "Legacy Source", feed.Items[0].Title)
	assert.Equal(t, "Elsewhere", feed.Items[1].Title)
	require.NotNil(t, feed.Items[1].Link)
	assert.Equal(t, "https://elsewhere.test/path", feed.Items[1].Link.Href)
}
