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
	var link string
	// Try to get href from the element itself
	if href, exists := sel.Attr("href"); exists {
		link = href
	} else {
		// If not found, try to find a child 'a' tag
		childA := sel.Find("a").First()
		if childA.Length() > 0 {
			if href, exists := childA.Attr("href"); exists {
				link = href
			}
		}

		// If still not found, try to find a parent 'a' tag (closest ancestor)
		if link == "" {
			parentA := sel.Closest("a")
			if parentA.Length() > 0 {
				if href, exists := parentA.Attr("href"); exists {
					link = href
				}
			}
		}
	}

	// Fallback to text if still empty, BUT only if it looks like a URL (no spaces)
	if link == "" {
		text := strings.TrimSpace(sel.Text())
		if text != "" && !strings.ContainsAny(text, " \t\n") {
			link = text
		}
	}
	return link
}
