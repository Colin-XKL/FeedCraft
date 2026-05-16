package examplefeed

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"FeedCraft/internal/model"

	"github.com/google/uuid"
)

const (
	RoutePrefix       = "/example-rss-feeds"
	CatalogAPIPath    = "/api/example-rss-feeds"
	AssetRoutePrefix  = RoutePrefix + "/assets"
	uuidWindow        = 4 * time.Hour
	uuidNamespaceSeed = "feedcraft-example-rss-feed"
)

var ErrUnknownFeed = errors.New("unknown example rss feed")

type FeedMeta struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

type feedDefinition struct {
	FeedMeta
	sections []contentSection
}

type contentSection struct {
	key         string
	title       string
	description string
	html        string
}

var feedDefinitions = []feedDefinition{
	{
		FeedMeta: FeedMeta{
			Slug:        "html-elements",
			Title:       "FeedCraft Example RSS Feeds - HTML Elements",
			Description: "Exercises common HTML5 elements in RSS item content.",
			Path:        RoutePrefix + "/html-elements.xml",
		},
		sections: []contentSection{htmlElementsSection},
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "html-styling",
			Title:       "FeedCraft Example RSS Feeds - HTML Styling",
			Description: "Exercises common inline CSS styles in RSS item content.",
			Path:        RoutePrefix + "/html-styling.xml",
		},
		sections: []contentSection{htmlStylingSection},
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "media-picture",
			Title:       "FeedCraft Example RSS Feeds - Picture Source",
			Description: "Exercises picture, source, srcset, image fallback, and captions.",
			Path:        RoutePrefix + "/media-picture.xml",
		},
		sections: []contentSection{mediaPictureSection},
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "all-in-one",
			Title:       "FeedCraft Example RSS Feeds - All in One",
			Description: "Combines HTML elements, inline styling, and media tests into one feed.",
			Path:        RoutePrefix + "/all-in-one.xml",
		},
		sections: []contentSection{htmlElementsSection, htmlStylingSection, mediaPictureSection},
	},
}

func Catalog() []FeedMeta {
	items := make([]FeedMeta, 0, len(feedDefinitions))
	for _, def := range feedDefinitions {
		items = append(items, def.FeedMeta)
	}
	return items
}

func Build(slug string, now time.Time) (*model.CraftFeed, error) {
	def, ok := findDefinition(slug)
	if !ok {
		return nil, ErrUnknownFeed
	}

	windowStart := now.UTC().Truncate(uuidWindow)
	rotatingID := WindowUUID(slug, now)
	articles := make([]*model.CraftArticle, 0, len(def.sections))
	for idx, section := range def.sections {
		articleID := fmt.Sprintf("%s:%s:%s", def.Slug, section.key, rotatingID)
		articles = append(articles, &model.CraftArticle{
			Title:       section.title,
			Link:        fmt.Sprintf("%s/%s.xml#%s", RoutePrefix, def.Slug, section.key),
			Description: section.description,
			Id:          articleID,
			Created:     windowStart,
			Updated:     windowStart,
			Content:     injectWindowID(section.html, rotatingID),
			AuthorName:  "FeedCraft",
		})
		if idx == 0 && len(def.sections) == 1 {
			articles[idx].Title = def.Title
		}
	}

	return &model.CraftFeed{
		Title:       def.Title,
		Link:        def.Path,
		Description: def.Description,
		Updated:     windowStart,
		Created:     windowStart,
		Id:          def.Path,
		AuthorName:  "FeedCraft",
		Articles:    articles,
	}, nil
}

func WindowUUID(slug string, now time.Time) string {
	windowStart := now.UTC().Truncate(uuidWindow)
	seed := fmt.Sprintf("%s:%s:%d", uuidNamespaceSeed, slug, windowStart.Unix())
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(seed)).String()
}

func findDefinition(slug string) (feedDefinition, bool) {
	for _, def := range feedDefinitions {
		if def.Slug == slug {
			return def, true
		}
	}
	return feedDefinition{}, false
}

func injectWindowID(html string, id string) string {
	return strings.ReplaceAll(html, "{{WINDOW_UUID}}", id)
}
