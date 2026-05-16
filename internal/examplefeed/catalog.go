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
	sections     []contentSection
	outputFormat outputFormat
}

type outputFormat string

const (
	outputRSS2     outputFormat = "rss2"
	outputRSS1     outputFormat = "rss1"
	outputRSS092   outputFormat = "rss0.92"
	outputAtom     outputFormat = "atom"
	outputJSONFeed outputFormat = "json-feed"
)

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
		sections:     []contentSection{htmlElementsSection},
		outputFormat: outputRSS2,
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "html-styling",
			Title:       "FeedCraft Example RSS Feeds - HTML Styling",
			Description: "Exercises common inline CSS styles in RSS item content.",
			Path:        RoutePrefix + "/html-styling.xml",
		},
		sections:     []contentSection{htmlStylingSection},
		outputFormat: outputRSS2,
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "media-picture",
			Title:       "FeedCraft Example RSS Feeds - Picture Source",
			Description: "Exercises picture, source, srcset, image fallback, and captions.",
			Path:        RoutePrefix + "/media-picture.xml",
		},
		sections:     []contentSection{mediaPictureSection},
		outputFormat: outputRSS2,
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "all-in-one",
			Title:       "FeedCraft Example RSS Feeds - All in One",
			Description: "Combines HTML elements, inline styling, and media tests into one feed.",
			Path:        RoutePrefix + "/all-in-one.xml",
		},
		sections:     []contentSection{htmlElementsSection, htmlStylingSection, mediaPictureSection},
		outputFormat: outputRSS2,
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "rss-2-0",
			Title:       "FeedCraft Example RSS Feeds - RSS 2.0",
			Description: "A simple RSS 2.0 feed for checking format support.",
			Path:        RoutePrefix + "/rss-2-0.xml",
		},
		sections:     []contentSection{formatSection},
		outputFormat: outputRSS2,
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "rss-1-0",
			Title:       "FeedCraft Example RSS Feeds - RSS 1.0",
			Description: "A simple RSS 1.0/RDF feed for checking format support.",
			Path:        RoutePrefix + "/rss-1-0.rdf",
		},
		sections:     []contentSection{formatSection},
		outputFormat: outputRSS1,
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "rss-0-92",
			Title:       "FeedCraft Example RSS Feeds - RSS 0.92",
			Description: "A simple legacy RSS 0.92 feed for checking reader compatibility.",
			Path:        RoutePrefix + "/rss-0-92.xml",
		},
		sections:     []contentSection{formatSection},
		outputFormat: outputRSS092,
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "atom",
			Title:       "FeedCraft Example RSS Feeds - Atom",
			Description: "A simple Atom feed for checking format support.",
			Path:        RoutePrefix + "/atom.xml",
		},
		sections:     []contentSection{formatSection},
		outputFormat: outputAtom,
	},
	{
		FeedMeta: FeedMeta{
			Slug:        "json-feed",
			Title:       "FeedCraft Example RSS Feeds - JSON Feed",
			Description: "A simple JSON Feed 1.1 document for checking format support.",
			Path:        RoutePrefix + "/json-feed.json",
		},
		sections:     []contentSection{formatSection},
		outputFormat: outputJSONFeed,
	},
}

func Catalog() []FeedMeta {
	items := make([]FeedMeta, 0, len(feedDefinitions))
	for _, def := range feedDefinitions {
		items = append(items, def.FeedMeta)
	}
	return items
}

func Build(slug string, now time.Time, baseURL string) (*model.CraftFeed, error) {
	def, ok := findDefinition(slug)
	if !ok {
		return nil, ErrUnknownFeed
	}

	baseURL = normalizeBaseURL(baseURL)
	feedURL := absoluteURL(baseURL, def.Path)
	windowStart := now.UTC().Truncate(uuidWindow)
	rotatingID := WindowUUID(slug, now)
	articles := make([]*model.CraftArticle, 0, len(def.sections))
	for idx, section := range def.sections {
		articleLink := fmt.Sprintf("%s#%s", feedURL, section.key)
		articles = append(articles, &model.CraftArticle{
			Title:       section.title,
			Link:        articleLink,
			Description: section.description,
			Id:          fmt.Sprintf("%s-%s", articleLink, rotatingID),
			Created:     windowStart,
			Updated:     windowStart,
			Content:     injectPlaceholders(section.html, rotatingID, baseURL),
			AuthorName:  "FeedCraft",
		})
		if idx == 0 && len(def.sections) == 1 {
			articles[idx].Title = def.Title
		}
	}

	return &model.CraftFeed{
		Title:       def.Title,
		Link:        feedURL,
		Description: def.Description,
		Updated:     windowStart,
		Created:     windowStart,
		Id:          feedURL,
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

func findDefinitionByPathName(pathName string) (feedDefinition, bool) {
	for _, def := range feedDefinitions {
		if def.Path == RoutePrefix+"/"+pathName {
			return def, true
		}
	}
	return feedDefinition{}, false
}

func injectPlaceholders(html string, id string, baseURL string) string {
	html = strings.ReplaceAll(html, "{{WINDOW_UUID}}", id)
	return strings.ReplaceAll(html, "{{BASE_URL}}", baseURL)
}

func absoluteURL(baseURL string, path string) string {
	if baseURL == "" {
		return path
	}
	return baseURL + path
}

func normalizeBaseURL(baseURL string) string {
	return strings.TrimRight(strings.TrimSpace(baseURL), "/")
}
