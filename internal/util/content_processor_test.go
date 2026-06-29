package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tinyBase64PNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="

func TestRemoveBase64Images_DropsDoubleQuotedImgTag(t *testing.T) {
	content := `<p>hello</p><img src="data:image/png;base64,` + tinyBase64PNG + `" alt="x"><p>world</p>`
	out := RemoveBase64Images(content)
	assert.NotContains(t, out, "base64")
	assert.NotContains(t, out, "<img")
	assert.Contains(t, out, "<p>hello</p>")
	assert.Contains(t, out, "<p>world</p>")
}

func TestRemoveBase64Images_DropsSingleQuotedImgTag(t *testing.T) {
	content := `<img alt='photo' src='data:image/jpeg;base64,` + tinyBase64PNG + `'/>keep`
	out := RemoveBase64Images(content)
	assert.NotContains(t, out, "base64")
	assert.NotContains(t, out, "<img")
	assert.Contains(t, out, "keep")
}

func TestRemoveBase64Images_DropsSourceTag(t *testing.T) {
	content := `<picture><source srcset="data:image/webp;base64,` + tinyBase64PNG + `"><img src="https://example.com/a.png"></picture>`
	out := RemoveBase64Images(content)
	assert.NotContains(t, out, "base64")
	assert.NotContains(t, out, "<source")
	// Regular (non-base64) image must be preserved.
	assert.Contains(t, out, `src="https://example.com/a.png"`)
}

func TestRemoveBase64Images_DropsMarkdownImage(t *testing.T) {
	content := "before ![alt](data:image/png;base64," + tinyBase64PNG + ") after"
	out := RemoveBase64Images(content)
	assert.NotContains(t, out, "base64")
	assert.Contains(t, out, "before")
	assert.Contains(t, out, "after")
}

func TestRemoveBase64Images_KeepsNormalImages(t *testing.T) {
	content := `<img src="https://example.com/pic.png" alt="ok">`
	out := RemoveBase64Images(content)
	assert.Equal(t, content, out)
}

func TestRemoveBase64Images_EmptyInput(t *testing.T) {
	assert.Equal(t, "", RemoveBase64Images(""))
}

func TestProcessContent_RemoveImageDropsBase64SingleQuote(t *testing.T) {
	content := `<p>text</p><img src='data:image/png;base64,` + tinyBase64PNG + `'>`
	out := ProcessContent(content, ContentProcessOption{RemoveImage: true})
	assert.NotContains(t, out, "base64")
	assert.NotContains(t, out, "<img")
	assert.Contains(t, out, "<p>text</p>")
}

func TestProcessContent_KeepsBase64WhenRemoveImageDisabled(t *testing.T) {
	content := `<img src="data:image/png;base64,` + tinyBase64PNG + `">`
	out := ProcessContent(content, ContentProcessOption{})
	assert.Equal(t, content, out)
}

func TestHTMLToMarkdown_StripsBase64Image(t *testing.T) {
	content := `<p>Article body here.</p><img src="data:image/png;base64,` + tinyBase64PNG + `" alt="diagram">`
	md := HTMLToMarkdown(content, "")
	assert.NotContains(t, md, "base64")
	assert.NotContains(t, md, "data:image")
	assert.Contains(t, md, "Article body here.")
	// Sanity check: the stripped output is far smaller than the base64 blob.
	assert.Less(t, len(md), len(tinyBase64PNG))
}

func TestHTMLToMarkdown_KeepsNormalImageLink(t *testing.T) {
	content := `<p>text</p><img src="https://example.com/pic.png" alt="ok">`
	md := HTMLToMarkdown(content, "")
	assert.Contains(t, strings.ToLower(md), "example.com/pic.png")
}
