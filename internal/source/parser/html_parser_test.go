package parser

import (
	"FeedCraft/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHtmlParserExtractsFeedIcon(t *testing.T) {
	htmlContent := `
		<html>
			<head>
				<title>Example Site</title>
				<link rel="stylesheet" href="/app.css">
				<link rel="shortcut icon" href="/favicon-32.png">
			</head>
			<body>
				<article class="item">
					<a class="title" href="/posts/1">First post</a>
				</article>
			</body>
		</html>`

	parser := &HtmlParser{Config: &config.HtmlParserConfig{
		ItemSelector: ".item",
		Title:        ".title",
		Link:         ".title",
	}}

	feed, err := parser.Parse([]byte(htmlContent))

	assert.NoError(t, err)
	if assert.NotNil(t, feed) {
		assert.Equal(t, "Example Site", feed.Title)
		assert.Equal(t, "/favicon-32.png", feed.ImageURL)
		assert.Equal(t, "Example Site", feed.ImageTitle)
	}
}
