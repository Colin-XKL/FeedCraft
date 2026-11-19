package util

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLLMContentCompactor_Base64MediaAndRestore(t *testing.T) {
	largePayload := strings.Repeat("QUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVo=", 4000) // ~144KB
	imgDataURI := "data:image/png;base64," + largePayload
	audioDataURI := "data:audio/mpeg;base64," + tinyBase64PNG

	input := fmt.Sprintf(
		`<p>Hello</p><img src="%s" alt="cover"><audio src='%s'></audio><p>World</p>`,
		imgDataURI,
		audioDataURI,
	)

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(input)

	assert.Equal(t, 2, compactor.Count())
	assert.NotContains(t, compacted, "base64")
	assert.NotContains(t, compacted, largePayload)
	assert.Less(t, len(compacted), 200)
	assert.Contains(t, compacted, `src="https://fc.local/m/1"`)
	assert.Contains(t, compacted, `src='https://fc.local/m/2'`)

	// Simulate LLM translating text while keeping placeholders untouched.
	llmOutput := strings.ReplaceAll(compacted, "Hello", "你好")
	llmOutput = strings.ReplaceAll(llmOutput, "World", "世界")

	restored := compactor.Restore(llmOutput)
	assert.Contains(t, restored, `<p>你好</p>`)
	assert.Contains(t, restored, `<p>世界</p>`)
	assert.Contains(t, restored, `src="`+imgDataURI+`"`)
	assert.Contains(t, restored, `src='`+audioDataURI+`'`)
}

func TestLLMContentCompactor_LongURLsAndShortURLs(t *testing.T) {
	shortURL := "https://example.com/posts/hello-world"
	longURL := "https://cdn.example.com/assets/image.png?" + strings.Repeat("sig=abcdef0123456789&", 10) + "end=1"
	require.GreaterOrEqual(t, len(longURL), DefaultLLMCompactorMaxURLLength)
	require.Less(t, len(shortURL), DefaultLLMCompactorMaxURLLength)

	input := fmt.Sprintf(
		`<p><a href="%s">short</a> and <a href="%s">long</a></p>`,
		shortURL,
		longURL,
	)

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(input)

	assert.Equal(t, 1, compactor.Count())
	assert.Contains(t, compacted, shortURL, "short URLs should remain intact")
	assert.NotContains(t, compacted, longURL, "long URLs should be compacted")
	assert.Contains(t, compacted, `href="https://fc.local/u/1"`)

	restored := compactor.Restore(compacted)
	assert.Equal(t, input, restored)
}

func TestLLMContentCompactor_MarkdownLinksWithParenthesesAndTrailingPunctuation(t *testing.T) {
	longWikiURL := "https://en.wikipedia.org/wiki/Some_Very_Long_Article_Title_(disambiguation)_" + strings.Repeat("x", 80)
	longPlainURL := "https://example.com/tracking?" + strings.Repeat("param=1234567890&", 10) + "ok=1"

	input := fmt.Sprintf(
		"Check ![pic](data:image/png;base64,%s) and [wiki](%s). Also see %s.",
		tinyBase64PNG,
		longWikiURL,
		longPlainURL,
	)

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(input)

	assert.Equal(t, 3, compactor.Count())
	assert.Contains(t, compacted, "![pic](https://fc.local/m/1)")
	assert.Contains(t, compacted, "[wiki](https://fc.local/u/1).")
	assert.Contains(t, compacted, "Also see https://fc.local/u/2.")

	assert.Equal(t, input, compactor.Restore(compacted))
}

func TestLLMContentCompactor_DeduplicatesIdenticalPayloads(t *testing.T) {
	dataURI := "data:image/png;base64," + tinyBase64PNG
	input := fmt.Sprintf(`<img src="%s"><img src="%s">`, dataURI, dataURI)

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(input)

	assert.Equal(t, 1, compactor.Count())
	assert.Equal(t, `<img src="https://fc.local/m/1"><img src="https://fc.local/m/1">`, compacted)
	assert.Equal(t, input, compactor.Restore(compacted))
}

func TestLLMContentCompactor_RestoreHandlesMoreThanTenPlaceholdersWithoutPrefixCollision(t *testing.T) {
	var parts []string
	urls := make([]string, 12)
	for i := 0; i < 12; i++ {
		urls[i] = fmt.Sprintf("https://example.com/item/%02d?token=%s", i, strings.Repeat("a", 120))
		parts = append(parts, fmt.Sprintf(`<a href="%s">link %d</a>`, urls[i], i))
	}
	input := strings.Join(parts, "\n")

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(input)

	assert.Equal(t, 12, compactor.Count())
	assert.Contains(t, compacted, "https://fc.local/u/1\"")
	assert.Contains(t, compacted, "https://fc.local/u/10\"")
	assert.Contains(t, compacted, "https://fc.local/u/12\"")

	restored := compactor.Restore(compacted)
	assert.Equal(t, input, restored)
}

func TestLLMContentCompactor_AvoidsPlaceholderCollisionWithExistingInput(t *testing.T) {
	longURL := "https://example.com/path?" + strings.Repeat("k=v&", 40)
	input := fmt.Sprintf("Mention https://fc.local/u/1 in text and link %s", longURL)

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(input)

	assert.Equal(t, 1, compactor.Count())
	assert.Contains(t, compacted, "Mention https://fc.local/u/1 in text")
	assert.Contains(t, compacted, "and link https://fc.local/u/2")
	assert.Equal(t, input, compactor.Restore(compacted))
}

func TestLLMContentCompactor_JSONEscapedAndMultilineBase64(t *testing.T) {
	jsonInput := fmt.Sprintf(`{"content":"<img src=\"data:image/png;base64,%s\">"}`, tinyBase64PNG)
	compactor := NewLLMContentCompactor()
	compactedJSON := compactor.Compact(jsonInput)
	assert.NotContains(t, compactedJSON, "base64")
	assert.Contains(t, compactedJSON, `\"https://fc.local/m/1\"`)
	assert.Equal(t, jsonInput, compactor.Restore(compactedJSON))

	multilineDataURI := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUg\nAAAAEAAAABCAQAAAC1HAwC"
	multilineHTML := `<img src="` + multilineDataURI + `">`
	compactor2 := NewLLMContentCompactor()
	compactedMultiline := compactor2.Compact(multilineHTML)
	assert.Equal(t, `<img src="https://fc.local/m/1">`, compactedMultiline)
	assert.Equal(t, multilineHTML, compactor2.Restore(compactedMultiline))
}

func TestLLMContentCompactor_SkipsActiveDataURIsAndRespectsWordBoundary(t *testing.T) {
	activeHTMLDataURI := "data:text/html;base64," + strings.Repeat("PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==", 10)
	activeSVGDataURI := "data:image/svg+xml;utf8," + strings.Repeat("%3Crect%20width%3D%2210%22/%3E", 10)
	safeFontDataURI := "data:font/woff2;base64," + tinyBase64PNG

	input := fmt.Sprintf(
		`metadata:image/png;base64,%s <a href="%s">bad</a> <img src="%s"> <link href="%s">`,
		tinyBase64PNG,
		activeHTMLDataURI,
		activeSVGDataURI,
		safeFontDataURI,
	)

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(input)

	assert.Equal(t, 1, compactor.Count())
	assert.Contains(t, compacted, "metadata:image/png;base64,"+tinyBase64PNG)
	assert.Contains(t, compacted, activeHTMLDataURI, "active data:text/html URI must not be compacted")
	assert.Contains(t, compacted, activeSVGDataURI, "active data:image/svg+xml URI must not be compacted")
	assert.Contains(t, compacted, `href="https://fc.local/m/1"`)
	assert.Equal(t, input, compactor.Restore(compacted))
}

func TestLLMContentCompactor_ReservePreventsPromptAndPrefixCollisions(t *testing.T) {
	longURL := "https://example.com/path?" + strings.Repeat("k=v&", 40)
	prompt := "Keep literal https://fc.local/u/12 untouched in output."
	input := "Read more at " + longURL

	compactor := NewLLMContentCompactor()
	compactor.Reserve(prompt)
	compacted := compactor.Compact(input)

	// https://fc.local/u/12 reserves both u/1 (prefix) and u/12.
	assert.Equal(t, "Read more at https://fc.local/u/2", compacted)

	llmOutput := "Prompt said https://fc.local/u/12 and link is https://fc.local/u/2"
	restored := compactor.Restore(llmOutput)
	assert.Contains(t, restored, "Prompt said https://fc.local/u/12")
	assert.Contains(t, restored, "link is "+longURL)
}

func TestLLMContentCompactor_JSONEscapedSlashes(t *testing.T) {
	escapedDataURI := `data:image\/png;base64,iVBORw0KGgo\/AAAANSUhEUg` + strings.Repeat(`A\/B`, 40)
	escapedLongURL := `https:\/\/cdn.example.com\/articles\/view?` + strings.Repeat(`a=1234567890&`, 12) + `end=1`
	jsonInput := fmt.Sprintf(`{"img":"%s","url":"%s"}`, escapedDataURI, escapedLongURL)

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(jsonInput)

	assert.Equal(t, 2, compactor.Count())
	assert.NotContains(t, compacted, "iVBORw0KGgo")
	assert.NotContains(t, compacted, "cdn.example.com")
	assert.Equal(t, `{"img":"https://fc.local/m/1","url":"https://fc.local/u/1"}`, compacted)
	assert.Equal(t, jsonInput, compactor.Restore(compacted))
}

func TestLLMContentCompactor_QuotedHTMLAttributeWithUnmatchedParenthesis(t *testing.T) {
	longURLWithParen := "https://example.com/path/foo)bar?" + strings.Repeat("token=abcdef0123&", 10) + "ok=1"
	input := fmt.Sprintf(`<a href="%s">link</a> and [md](%s)`, longURLWithParen, "https://example.com/md?"+strings.Repeat("x=1234567890&", 10))

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(input)

	assert.Equal(t, 2, compactor.Count())
	assert.Equal(t, `<a href="https://fc.local/u/1">link</a> and [md](https://fc.local/u/2)`, compacted)
	assert.Equal(t, input, compactor.Restore(compacted))
}

func TestLLMContentCompactor_SurvivesHTMLMarkdownRoundTrip(t *testing.T) {
	dataURI := "data:image/png;base64," + tinyBase64PNG
	longURL := "https://example.com/article?a=1&amp;b=2&amp;token=" + strings.Repeat("z", 120)
	htmlInput := fmt.Sprintf(
		`<p>Intro text</p><p><img src="%s" alt="diagram"></p><p><a href="%s">Read more</a></p>`,
		dataURI,
		longURL,
	)

	compactor := NewLLMContentCompactor()
	compactedHTML := compactor.Compact(htmlInput)

	md := HTMLToMarkdown(compactedHTML, "https://blog.example.org")
	assert.Contains(t, md, "https://fc.local/m/1", "placeholder must not be rewritten with domain")
	assert.Contains(t, md, "https://fc.local/u/1", "placeholder must not be rewritten with domain")
	assert.NotContains(t, md, "blog.example.org/https://fc.local")

	renderedHTML := MarkdownToHTML(md)
	restoredHTML := compactor.Restore(renderedHTML)
	assert.Contains(t, restoredHTML, `src="`+dataURI+`"`)
	assert.Contains(t, restoredHTML, `href="`+longURL+`"`)
}

func TestLLMContentCompactor_FourMegabyteBase64Image(t *testing.T) {
	fourMBBase64 := "data:image/png;base64," + strings.Repeat("A", 4*1024*1024)
	input := `<p>Title</p><img src="` + fourMBBase64 + `" alt="huge"><p>Footer</p>`

	compactor := NewLLMContentCompactor()
	compacted := compactor.Compact(input)

	assert.Equal(t, 1, compactor.Count())
	assert.Less(t, len(compacted), 100)
	assert.Equal(t, `<p>Title</p><img src="https://fc.local/m/1" alt="huge"><p>Footer</p>`, compacted)

	restored := compactor.Restore(compacted)
	assert.Equal(t, input, restored)
}
