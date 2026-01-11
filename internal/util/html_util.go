package util

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractLinkFromSelection tries to extract a link URL from a goquery Selection.
// It uses a heuristic approach:
// 1. Checks 'href' attribute on the selection itself.
// 2. Checks 'href' attribute on the first child 'a' tag.
// 3. Checks 'href' attribute on the closest parent 'a' tag.
// 4. Fallbacks to text content if it looks like a URL (no spaces).
func ExtractLinkFromSelection(sel *goquery.Selection) string {
	// 1. Try to get href from the element itself
	if href, exists := sel.Attr("href"); exists {
		return href
	}

	// 2. If not found, try to find a child 'a' tag
	if href, exists := sel.Find("a").First().Attr("href"); exists {
		return href
	}

	// 3. If still not found, try to find a parent 'a' tag (closest ancestor)
	if href, exists := sel.Closest("a").Attr("href"); exists {
		return href
	}

	// 4. Fallback to text if still empty, BUT only if it looks like a URL (no spaces)
	text := strings.TrimSpace(sel.Text())
	if text != "" && !strings.ContainsAny(text, " \t\n") {
		return text
	}

	return ""
}
