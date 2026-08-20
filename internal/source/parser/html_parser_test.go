package parser

import (
	"FeedCraft/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestHtmlParserIgnoresUnrelatedRelContainingIconSubstring(t *testing.T) {
	htmlContent := `
		<html>
			<head>
				<title>Example Site</title>
				<link rel="iconography" href="/not-a-favicon.png">
				<link rel="preload" href="/icon.svg">
				<link rel="apple-touch-icon" href="/apple-touch.png">
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

	require.NoError(t, err)
	require.NotNil(t, feed)
	assert.Equal(t, "/apple-touch.png", feed.ImageURL)
}

func TestHasFeedIconRel(t *testing.T) {
	assert.True(t, hasFeedIconRel("icon"))
	assert.True(t, hasFeedIconRel("SHORTCUT ICON"))
	assert.True(t, hasFeedIconRel("apple-touch-icon"))
	assert.True(t, hasFeedIconRel("apple-touch-icon-precomposed"))
	assert.True(t, hasFeedIconRel("mask-icon"))
	assert.False(t, hasFeedIconRel("iconography"))
	assert.False(t, hasFeedIconRel("stylesheet"))
	assert.False(t, hasFeedIconRel("preload"))
	assert.False(t, hasFeedIconRel(""))
}
