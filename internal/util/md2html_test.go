package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdown2HTML_SingleNewlineBecomesBreak(t *testing.T) {
	html := Markdown2HTML("line1\nline2")
	assert.Contains(t, html, "line1<br>")
	assert.Contains(t, html, "line2")
}

func TestMarkdown2HTML_DoubleNewlineCreatesParagraphs(t *testing.T) {
	html := Markdown2HTML("line1\n\nline2")
	assert.Contains(t, html, "<p>line1</p>")
	assert.Contains(t, html, "<p>line2</p>")
}

func TestMarkdown2HTML_TwoSpacesNewlineStillBreaks(t *testing.T) {
	html := Markdown2HTML("line1  \nline2")
	assert.Contains(t, html, "line1<br>")
	assert.Contains(t, html, "line2")
}

func TestMarkdown2HTML_ExplicitBreakTagPreserved(t *testing.T) {
	html := Markdown2HTML("line1<br>line2")
	assert.Contains(t, html, "line1<br>")
	assert.Contains(t, html, "line2")
}

func TestMarkdown2HTML_ListItemsUnchanged(t *testing.T) {
	html := Markdown2HTML("- item1\n- item2")
	assert.Contains(t, html, "<li>item1</li>")
	assert.Contains(t, html, "<li>item2</li>")
	assert.NotRegexp(t, `<li>item1<br>`, html)
}

func TestMarkdown2HTML_BlockquotePreservesLineBreaks(t *testing.T) {
	html := Markdown2HTML("> quote line1\n> quote line2")
	assert.Contains(t, html, "quote line1<br>")
	assert.Contains(t, html, "quote line2")
}

func TestMarkdown2HTML_CodeBlockPreservesNewlines(t *testing.T) {
	html := Markdown2HTML("```\ncode line1\ncode line2\n```")
	assert.Contains(t, html, "<pre>")
	assert.Contains(t, html, "code line1")
	assert.Contains(t, html, "code line2")
	assert.NotContains(t, html, "code line1<br>")
}

func TestInsertLineBreaksInParagraphHTML_SkipsLeadingAndTrailingNewlines(t *testing.T) {
	html := insertLineBreaksInParagraphHTML("<p>\nline1\nline2\n</p>")
	assert.Contains(t, html, "line1<br>")
	assert.Contains(t, html, "line2")
	assert.NotContains(t, html, "<br>\n\n")
}

func TestInsertLineBreaksInParagraphHTML_DoesNotDuplicateExistingBreak(t *testing.T) {
	html := insertLineBreaksInParagraphHTML("<p>line1<br>\nline2</p>")
	assert.Equal(t, "<p>line1<br>\nline2</p>", html)
}

func TestInsertLineBreaksInParagraphHTML_HandlesAttributes(t *testing.T) {
	html := insertLineBreaksInParagraphHTML(`<p class="note">line1
line2</p>`)
	assert.Contains(t, html, `class="note"`)
	assert.Contains(t, html, "line1<br>")
}

func TestHtml2Markdown_ParagraphNewlineRoundTrip(t *testing.T) {
	domain := "example.com"
	original := "<p>line1\nline2</p>"

	md := Html2Markdown(original, &domain)
	require.NotEmpty(t, md)
	assert.Contains(t, md, "line1")
	assert.Contains(t, md, "line2")

	back := Markdown2HTML(md)
	assert.Contains(t, back, "line1<br>")
	assert.Contains(t, back, "line2")
}

func TestHtml2Markdown_BreakTagRoundTrip(t *testing.T) {
	domain := "example.com"
	original := "<p>line1<br>line2</p>"

	md := Html2Markdown(original, &domain)
	back := Markdown2HTML(md)

	assert.Contains(t, back, "line1<br>")
	assert.Contains(t, back, "line2")
}

func TestHtml2Markdown_SeparateParagraphsRoundTrip(t *testing.T) {
	domain := "example.com"
	original := "<p>line1</p><p>line2</p>"

	md := Html2Markdown(original, &domain)
	back := Markdown2HTML(md)

	assert.Contains(t, back, "<p>line1</p>")
	assert.Contains(t, back, "<p>line2</p>")
}

func TestMarkdown2HTML_TableUnaffected(t *testing.T) {
	md := "| a | b |\n|---|---|\n| 1 | 2 |"
	html := Markdown2HTML(md)
	assert.Contains(t, html, "<table>")
	assert.Contains(t, html, ">1<")
	assert.Contains(t, html, ">2<")
}

func TestShouldInsertLineBreak(t *testing.T) {
	assert.True(t, shouldInsertLineBreak("line1", "line2"))
	assert.False(t, shouldInsertLineBreak("", "line2"))
	assert.False(t, shouldInsertLineBreak("line1", ""))
	assert.False(t, shouldInsertLineBreak("line1<br>", "line2"))
	assert.False(t, shouldInsertLineBreak("line1<br/>", "line2"))
}

func TestMarkdown2HTML_MultiLineParagraph(t *testing.T) {
	md := "first line\nsecond line\n\nnew paragraph"
	html := Markdown2HTML(md)

	firstP := strings.Index(html, "<p>")
	secondP := strings.Index(html[firstP+len("<p>"):], "<p>")
	require.NotEqual(t, -1, secondP)
	secondP += firstP + len("<p>")
	require.NotEqual(t, -1, firstP)
	require.NotEqual(t, -1, secondP)

	firstParagraph := html[firstP:secondP]
	assert.Contains(t, firstParagraph, "first line<br>")
	assert.Contains(t, firstParagraph, "second line")
}
