package util

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestExtractLinkFromSelection(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		selector string
		expected string
	}{
		{
			name:     "Link in self href",
			html:     `<a href="https://example.com/1">Link 1</a>`,
			selector: "a",
			expected: "https://example.com/1",
		},
		{
			name:     "Link in child a",
			html:     `<div><a href="https://example.com/2">Link 2</a></div>`,
			selector: "div",
			expected: "https://example.com/2",
		},
		{
			name:     "Link in parent a",
			html:     `<a href="https://example.com/3"><span>Link 3</span></a>`,
			selector: "span",
			expected: "https://example.com/3",
		},
		{
			name:     "Link in text",
			html:     `<div>https://example.com/4</div>`,
			selector: "div",
			expected: "https://example.com/4",
		},
		{
			name:     "Link in text with whitespace",
			html:     `<div>  https://example.com/5  </div>`,
			selector: "div",
			expected: "https://example.com/5",
		},
		{
			name:     "No link (text with spaces)",
			html:     `<div>Not a link</div>`,
			selector: "div",
			expected: "",
		},
		{
			name:     "No link (empty)",
			html:     `<div></div>`,
			selector: "div",
			expected: "",
		},
		{
			name:     "Priority: Self over Child",
			html:     `<a href="https://self.com"><span><a href="https://child.com">Child</a></span></a>`,
			selector: "a",
			expected: "https://self.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			assert.NoError(t, err)
			sel := doc.Find(tt.selector).First()
			got := ExtractLinkFromSelection(sel)
			assert.Equal(t, tt.expected, got)
		})
	}
}
