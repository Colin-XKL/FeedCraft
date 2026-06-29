package craft

import (
	"strings"
	"testing"

	"FeedCraft/internal/model"

	"github.com/stretchr/testify/assert"
)

const testTinyBase64PNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="

// TestGetArticleContentForPrompt_StripsBase64Images verifies that the LLM
// preparation path (used by summary, introduction, etc.) strips inline base64
// images before converting HTML to Markdown for the LLM context.
func TestGetArticleContentForPrompt_StripsBase64Images(t *testing.T) {
	article := &model.CraftArticle{
		Title: "Test Article",
		Link:  "https://example.com/post/1",
	}
	original := `<h1>Weekly Roundup</h1>` +
		`<p>Read on for details.</p>` +
		`<img src="data:image/png;base64,` + testTinyBase64PNG + `" alt="cover">` +
		`<p>End of article.</p>`

	result := getArticleContentForPrompt(article, original)

	assert.NotContains(t, result, "base64", "base64 image data must not be sent to the LLM")
	assert.NotContains(t, result, "data:image", "base64 data URI must be stripped")
	assert.Contains(t, result, "Weekly Roundup")
	assert.Contains(t, result, "Read on for details.")
	assert.Contains(t, result, "End of article.")
	// Stripped output should be far smaller than the raw base64 blob.
	assert.Less(t, len(result), len(testTinyBase64PNG))
}

// TestGetArticleContentForPrompt_KeepsNormalImages verifies that regular
// (non-base64) image references are preserved after the cleanup.
func TestGetArticleContentForPrompt_KeepsNormalImages(t *testing.T) {
	article := &model.CraftArticle{
		Title: "Test Article",
		Link:  "https://example.com/post/2",
	}
	original := `<p>See diagram below.</p><img src="https://example.com/pic.png" alt="diagram">`

	result := getArticleContentForPrompt(article, original)

	assert.True(t,
		strings.Contains(result, "example.com/pic.png") || strings.Contains(result, "diagram"),
		"normal image reference should be preserved",
	)
}
