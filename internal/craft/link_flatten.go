package craft

import (
	"context"
	"fmt"
	"html"
	"net/url"
	"strconv"
	"strings"

	"FeedCraft/internal/engine"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

type LinkFlattenConfig struct {
	KeepOriginal bool
	ExternalOnly bool
}

type linkCandidate struct {
	Title string
	URL   string
}

var linkFlattenParamTmpl = []ParamTemplate{
	{
		Key:         "keep-original",
		Description: "`true` to keep original RSS articles before extracted link items, `false` to output extracted links only",
		Default:     "true",
	},
	{
		Key:         "external-only",
		Description: "`true` to extract only links whose domain differs from the source article URL, `false` to include same-domain links too",
		Default:     "true",
	},
}

func linkFlattenCraftLoadParam(m map[string]string) []LegacyCraftOption {
	return GetLinkFlattenCraftOptions(parseLinkFlattenConfig(m))
}

func GetLinkFlattenCraftOptions(config LinkFlattenConfig) []LegacyCraftOption {
	return []LegacyCraftOption{OptionLinkFlatten(config)}
}

func OptionLinkFlatten(config LinkFlattenConfig) LegacyCraftOption {
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		if feed == nil || len(feed.Items) == 0 {
			return nil
		}

		feedLink := ""
		if feed.Link != nil {
			feedLink = feed.Link.Href
		}

		flattened := make([]*feeds.Item, 0, len(feed.Items))
		for _, item := range feed.Items {
			if item == nil {
				continue
			}
			if config.KeepOriginal {
				flattened = append(flattened, item)
			}

			sourceLink := ""
			if item.Link != nil {
				sourceLink = item.Link.Href
			}
			sourceLink = resolveArticleLink(payload.originalFeedUrl, feedLink, sourceLink)
			for _, candidate := range extractFlattenedLinks(item.Content, item.Description, sourceLink, config.ExternalOnly) {
				description := buildFlattenedLinkDescription(item.Title, sourceLink)
				flattened = append(flattened, &feeds.Item{
					Title:       candidate.Title,
					Link:        &feeds.Link{Href: candidate.URL},
					Description: description,
					Content:     description,
					Id:          buildFlattenedLinkID(item.Id, sourceLink, candidate.URL),
					Updated:     item.Updated,
					Created:     item.Created,
				})
			}
		}

		feed.Items = flattened
		return nil
	}
}

type linkFlattenProcessor struct {
	feedURL string
	config  LinkFlattenConfig
}

func newLinkFlattenProcessor(feedURL string, config LinkFlattenConfig) *linkFlattenProcessor {
	return &linkFlattenProcessor{feedURL: feedURL, config: config}
}

func (p *linkFlattenProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	_ = ctx
	if feed == nil || len(feed.Articles) == 0 {
		return feed, nil
	}

	cloned := cloneCraftFeed(feed)
	flattened := make([]*model.CraftArticle, 0, len(cloned.Articles))
	for _, article := range cloned.Articles {
		if article == nil {
			continue
		}
		if p.config.KeepOriginal {
			flattened = append(flattened, article)
		}

		sourceLink := resolveArticleLink(p.feedURL, cloned.Link, article.Link)
		for _, candidate := range extractFlattenedLinks(article.Content, article.Description, sourceLink, p.config.ExternalOnly) {
			flattened = append(flattened, &model.CraftArticle{
				Title:       candidate.Title,
				Link:        candidate.URL,
				Description: buildFlattenedLinkDescription(article.Title, sourceLink),
				Content:     buildFlattenedLinkDescription(article.Title, sourceLink),
				Id:          buildFlattenedLinkID(article.Id, sourceLink, candidate.URL),
				Updated:     article.Updated,
				Created:     article.Created,
				AuthorName:  article.AuthorName,
				AuthorEmail: article.AuthorEmail,
			})
		}
	}
	cloned.Articles = flattened
	return cloned, nil
}

func buildNativeLinkFlattenProcessor(atom ResolvedCraftAtom, feedURL string) (engine.CraftOption, error) {
	return wrapLocalProcessor(newLinkFlattenProcessor(feedURL, parseLinkFlattenConfig(atom.Params))), nil
}

func parseLinkFlattenConfig(params map[string]string) LinkFlattenConfig {
	return LinkFlattenConfig{
		KeepOriginal: parseBoolParamWithDefault(params["keep-original"], true, "keep-original"),
		ExternalOnly: parseBoolParamWithDefault(params["external-only"], true, "external-only"),
	}
}

func parseBoolParamWithDefault(raw string, defaultValue bool, paramName string) bool {
	if strings.TrimSpace(raw) == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(raw))
	if err != nil {
		logrus.Warnf("invalid bool param [%s]=[%s], using default [%v]", paramName, raw, defaultValue)
		return defaultValue
	}
	return parsed
}

func extractFlattenedLinks(content, description, sourceArticleURL string, externalOnly bool) []linkCandidate {
	htmlContent := content
	if strings.TrimSpace(description) != "" && description != content {
		htmlContent += description
	}
	if strings.TrimSpace(htmlContent) == "" {
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		logrus.Warnf("failed to parse article html for link flatten: %v", err)
		return nil
	}

	seen := make(map[string]struct{})
	candidates := make([]linkCandidate, 0)
	doc.Find("a[href]").Each(func(i int, selection *goquery.Selection) {
		rawHref, exists := selection.Attr("href")
		if !exists {
			return
		}

		absURL, ok := normalizeHTTPURL(sourceArticleURL, rawHref)
		if !ok {
			return
		}
		if externalOnly && !isExternalLink(sourceArticleURL, absURL) {
			return
		}
		if _, exists := seen[absURL]; exists {
			return
		}
		seen[absURL] = struct{}{}

		title := strings.TrimSpace(selection.Text())
		if title == "" {
			title = absURL
		}
		candidates = append(candidates, linkCandidate{Title: title, URL: absURL})
	})
	return candidates
}

func normalizeHTTPURL(baseURL, rawHref string) (string, bool) {
	trimmed := strings.TrimSpace(rawHref)
	if trimmed == "" {
		return "", false
	}

	parsedHref, err := url.Parse(trimmed)
	if err != nil {
		return "", false
	}

	var resolved *url.URL
	if parsedHref.IsAbs() {
		resolved = parsedHref
	} else {
		parsedBase, err := url.Parse(baseURL)
		if err != nil || parsedBase == nil || !parsedBase.IsAbs() {
			return "", false
		}
		resolved = parsedBase.ResolveReference(parsedHref)
	}

	if resolved == nil || (resolved.Scheme != "http" && resolved.Scheme != "https") || resolved.Hostname() == "" {
		return "", false
	}
	return resolved.String(), true
}

func isExternalLink(sourceArticleURL, targetURL string) bool {
	source, err := url.Parse(sourceArticleURL)
	if err != nil || source == nil {
		return false
	}
	target, err := url.Parse(targetURL)
	if err != nil || target == nil {
		return false
	}
	return !strings.EqualFold(source.Hostname(), target.Hostname())
}

func resolveArticleLink(feedURL, feedLink, articleLink string) string {
	if strings.TrimSpace(articleLink) != "" {
		return getAbsLinkForFeedItem(feedURL, feedLink, articleLink)
	}
	if strings.TrimSpace(feedLink) != "" {
		return getAbsFeedLink(feedURL, feedLink)
	}
	return feedURL
}

func buildFlattenedLinkDescription(sourceTitle, sourceURL string) string {
	return fmt.Sprintf(
		`Extracted from <a href="%s">%s</a>`,
		html.EscapeString(sourceURL),
		html.EscapeString(sourceTitle),
	)
}

func buildFlattenedLinkID(sourceID, sourceURL, targetURL string) string {
	return util.GetTextContentHash(sourceID + "|" + sourceURL + "|" + targetURL)
}
