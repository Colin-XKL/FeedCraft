package model

import (
	"time"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
)

// CraftFeed is the standard internal data carrier for the FeedCraft pipeline.
type CraftFeed struct {
	Title       string
	Link        string // The actual URL of the site
	Description string
	Updated     time.Time
	Created     time.Time
	Id          string // The URL of the feed itself
	Copyright   string
	AuthorName  string
	AuthorEmail string
	ImageURL    string
	ImageTitle  string

	// The items in the feed
	Articles []*CraftArticle
}

// CraftArticle represents a single item in the feed pipeline.
// It contains standard RSS fields as well as internal tracking metadata.
type CraftArticle struct {
	Title       string
	Link        string
	Description string
	Id          string // GUID
	Updated     time.Time
	Created     time.Time
	Content     string
	AuthorName  string
	AuthorEmail string

	// Internal metadata for future expansion (Topic aggregation, fission, etc.)
	OriginalFeedID string // The ID of the raw source it came from
	QualityScore   float64
	Depth          int // For tracking fission loops
}

// FromGofeed converts an external gofeed.Feed (usually from the Source layer) into a CraftFeed.
func FromGofeed(parsedFeed *gofeed.Feed) *CraftFeed {
	if parsedFeed == nil {
		return nil
	}

	updatedTime := time.Now()
	if parsedFeed.UpdatedParsed != nil && !parsedFeed.UpdatedParsed.IsZero() {
		updatedTime = *parsedFeed.UpdatedParsed
	}

	publishedTime := time.Now()
	if parsedFeed.PublishedParsed != nil && !parsedFeed.PublishedParsed.IsZero() {
		publishedTime = *parsedFeed.PublishedParsed
	}

	cf := &CraftFeed{
		Title:       parsedFeed.Title,
		Link:        parsedFeed.Link,
		Description: parsedFeed.Description,
		Updated:     updatedTime,
		Created:     publishedTime,
		Id:          parsedFeed.FeedLink,
		Copyright:   parsedFeed.Copyright,
	}

	if len(parsedFeed.Authors) > 0 && parsedFeed.Authors[0] != nil {
		cf.AuthorName = parsedFeed.Authors[0].Name
		cf.AuthorEmail = parsedFeed.Authors[0].Email
	} else if parsedFeed.Author != nil {
		cf.AuthorName = parsedFeed.Author.Name
		cf.AuthorEmail = parsedFeed.Author.Email
	}
	if parsedFeed.Image != nil {
		cf.ImageURL = parsedFeed.Image.URL
		cf.ImageTitle = parsedFeed.Image.Title
	}

	cf.Articles = lo.Map(parsedFeed.Items, func(item *gofeed.Item, _ int) *CraftArticle {
		return ArticleFromGofeed(item)
	})

	return cf
}

func ArticleFromGofeed(item *gofeed.Item) *CraftArticle {
	if item == nil {
		return nil
	}

	updatedTime := time.Now()
	if item.UpdatedParsed != nil && !item.UpdatedParsed.IsZero() {
		updatedTime = *item.UpdatedParsed
	}

	publishedTime := time.Now()
	if item.PublishedParsed != nil && !item.PublishedParsed.IsZero() {
		publishedTime = *item.PublishedParsed
	}

	content := item.Content
	if len(content) == 0 {
		content = item.Description
	}

	article := &CraftArticle{
		Title:       item.Title,
		Link:        item.Link,
		Description: item.Description,
		Id:          item.GUID,
		Updated:     updatedTime,
		Created:     publishedTime,
		Content:     content, // fallback to description if content is empty
	}

	authorItem := lo.FirstOrEmpty(item.Authors)
	if authorItem != nil {
		article.AuthorName = authorItem.Name
		article.AuthorEmail = authorItem.Email
	} else if item.Author != nil {
		article.AuthorName = item.Author.Name
		article.AuthorEmail = item.Author.Email
	}

	return article
}

// ToFeedsFeed converts the internal CraftFeed out to the Gorilla feeds format for final XML generation.
func (cf *CraftFeed) ToFeedsFeed() *feeds.Feed {
	if cf == nil {
		return nil
	}

	ff := &feeds.Feed{
		Title:       cf.Title,
		Link:        &feeds.Link{Href: cf.Link},
		Description: cf.Description,
		Updated:     cf.Updated,
		Created:     cf.Created,
		Id:          cf.Id,
		Copyright:   cf.Copyright,
	}

	if cf.AuthorName != "" || cf.AuthorEmail != "" {
		ff.Author = &feeds.Author{
			Name:  cf.AuthorName,
			Email: cf.AuthorEmail,
		}
	}
	if cf.ImageURL != "" {
		ff.Image = &feeds.Image{
			Url:   cf.ImageURL,
			Title: cf.ImageTitle,
			Link:  cf.ImageURL,
		}
	}

	ff.Items = lo.Map(cf.Articles, func(article *CraftArticle, _ int) *feeds.Item {
		return article.ToFeedsItem()
	})

	return ff
}

func (ca *CraftArticle) ToFeedsItem() *feeds.Item {
	if ca == nil {
		return nil
	}

	item := &feeds.Item{
		Title:       ca.Title,
		Link:        &feeds.Link{Href: ca.Link},
		Description: ca.Description,
		Id:          ca.Id,
		Updated:     ca.Updated,
		Created:     ca.Created,
		Content:     ca.Content,
	}

	if ca.AuthorName != "" || ca.AuthorEmail != "" {
		item.Author = &feeds.Author{
			Name:  ca.AuthorName,
			Email: ca.AuthorEmail,
		}
	}

	return item
}

// FromFeedsFeed converts a Gorilla feeds.Feed back to a CraftFeed.
// This is primarily used by the LegacyOptionAdapter to capture state after legacy CraftOptions run.
func FromFeedsFeed(ff *feeds.Feed) *CraftFeed {
	if ff == nil {
		return nil
	}

	cf := &CraftFeed{
		Title:       ff.Title,
		Description: ff.Description,
		Updated:     ff.Updated,
		Created:     ff.Created,
		Id:          ff.Id,
		Copyright:   ff.Copyright,
	}

	if ff.Link != nil {
		cf.Link = ff.Link.Href
	}
	if ff.Author != nil {
		cf.AuthorName = ff.Author.Name
		cf.AuthorEmail = ff.Author.Email
	}
	if ff.Image != nil {
		cf.ImageURL = ff.Image.Url
		cf.ImageTitle = ff.Image.Title
	}

	cf.Articles = lo.Map(ff.Items, func(item *feeds.Item, _ int) *CraftArticle {
		ca := &CraftArticle{
			Title:       item.Title,
			Description: item.Description,
			Id:          item.Id,
			Updated:     item.Updated,
			Created:     item.Created,
			Content:     item.Content,
		}
		if item.Link != nil {
			ca.Link = item.Link.Href
		}
		if item.Author != nil {
			ca.AuthorName = item.Author.Name
			ca.AuthorEmail = item.Author.Email
		}
		return ca
	})

	return cf
}
