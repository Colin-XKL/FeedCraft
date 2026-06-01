package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownToHTML_SingleNewlineBecomesBreak(t *testing.T) {
	html := MarkdownToHTML("line1\nline2")
	assert.Contains(t, html, "line1<br>")
	assert.Contains(t, html, "line2")
}

func TestMarkdownToHTML_DoubleNewlineCreatesParagraphs(t *testing.T) {
	html := MarkdownToHTML("line1\n\nline2")
	assert.Contains(t, html, "<p>line1</p>")
	assert.Contains(t, html, "<p>line2</p>")
}

func TestMarkdownToHTML_TwoSpacesNewlineStillBreaks(t *testing.T) {
	html := MarkdownToHTML("line1  \nline2")
	assert.Contains(t, html, "line1<br>")
	assert.Contains(t, html, "line2")
}

func TestMarkdownToHTML_ExplicitBreakTagPreserved(t *testing.T) {
	html := MarkdownToHTML("line1<br>line2")
	assert.Contains(t, html, "line1<br>")
	assert.Contains(t, html, "line2")
}

func TestMarkdownToHTML_ListItemsUnchanged(t *testing.T) {
	html := MarkdownToHTML("- item1\n- item2")
	assert.Contains(t, html, "<li>item1</li>")
	assert.Contains(t, html, "<li>item2</li>")
	assert.NotRegexp(t, `<li>item1<br>`, html)
}

func TestMarkdownToHTML_BlockquotePreservesLineBreaks(t *testing.T) {
	html := MarkdownToHTML("> quote line1\n> quote line2")
	assert.Contains(t, html, "quote line1<br>")
	assert.Contains(t, html, "quote line2")
}

func TestMarkdownToHTML_CodeBlockPreservesNewlines(t *testing.T) {
	html := MarkdownToHTML("```\ncode line1\ncode line2\n```")
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

func TestHTMLToMarkdown_ParagraphNewlineRoundTrip(t *testing.T) {
	original := "<p>line1\nline2</p>"

	md := HTMLToMarkdown(original, "example.com")
	require.NotEmpty(t, md)
	assert.Contains(t, md, "line1")
	assert.Contains(t, md, "line2")

	back := MarkdownToHTML(md)
	assert.Contains(t, back, "line1<br>")
	assert.Contains(t, back, "line2")
}

func TestHTMLToMarkdown_BreakTagRoundTrip(t *testing.T) {
	original := "<p>line1<br>line2</p>"

	md := HTMLToMarkdown(original, "example.com")
	back := MarkdownToHTML(md)

	assert.Contains(t, back, "line1<br>")
	assert.Contains(t, back, "line2")
}

func TestHTMLToMarkdown_SeparateParagraphsRoundTrip(t *testing.T) {
	original := "<p>line1</p><p>line2</p>"

	md := HTMLToMarkdown(original, "example.com")
	back := MarkdownToHTML(md)

	assert.Contains(t, back, "<p>line1</p>")
	assert.Contains(t, back, "<p>line2</p>")
}

func TestMarkdownToHTML_TableUnaffected(t *testing.T) {
	md := "| a | b |\n|---|---|\n| 1 | 2 |"
	html := MarkdownToHTML(md)
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

func TestCollapseConsecutiveBlankLines(t *testing.T) {
	assert.Equal(t, "line1\n\nline2", collapseConsecutiveBlankLines("line1\n\n\n\nline2"))
	assert.Equal(t, "line1\n\nline2", collapseConsecutiveBlankLines("line1\n   \n  \nline2"))
}

func TestApplySingleNewlineHardBreaks_ListBoundaryUnchanged(t *testing.T) {
	md := applySingleNewlineHardBreaks("- item1\n- item2")
	assert.Equal(t, "- item1\n- item2", md)
}

func TestApplySingleNewlineHardBreaks_BlockquoteContinuationGetsHardBreak(t *testing.T) {
	md := applySingleNewlineHardBreaks("> quote line1\n> quote line2")
	assert.Contains(t, md, "quote line1  \n> quote line2")
}

func TestApplySingleNewlineHardBreaks_ParagraphGetsHardBreak(t *testing.T) {
	md := applySingleNewlineHardBreaks("line1\nline2")
	assert.Equal(t, "line1  \nline2", md)
}

func TestNormalizeMarkdownForRender_PreservesCodeFenceBlankLines(t *testing.T) {
	input := "intro\n\n\n\n```\nkeep\n\n\nspacing\n```\n\n\n\noutro"
	normalized := normalizeMarkdownForRender(input)
	assert.Contains(t, normalized, "keep\n\n\nspacing")
	assert.Contains(t, normalized, "intro\n\n```")
	assert.Contains(t, normalized, "```\n\noutro")
}

func TestMarkdownToHTML_CollapsesExcessiveBlankLines(t *testing.T) {
	html := MarkdownToHTML("line1\n\n\n\nline2")
	assert.Contains(t, html, "<p>line1</p>")
	assert.Contains(t, html, "<p>line2</p>")
	assert.NotContains(t, html, "<p></p>")
}

func TestHTMLToMarkdown_CollapsesExcessiveBlankLines(t *testing.T) {
	md := HTMLToMarkdown("<p>line1</p><p>line2</p>", "example.com")
	assert.NotContains(t, md, "\n\n\n")
}

func TestMarkdownToHTML_MultiLineParagraph(t *testing.T) {
	md := "first line\nsecond line\n\nnew paragraph"
	html := MarkdownToHTML(md)

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
