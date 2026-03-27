package util

import (
	"regexp"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

var (
	linkRegex = regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"[^>]*>.*?</a>`)
	imgRegex  = regexp.MustCompile(`<img\s+(?:[^>]*?\s+)?src="([^"]*)"[^>]*>`)
)

type ContentProcessOption struct {
	RemoveLinks bool
	RemoveImage bool
	ConvertToMd bool
}

func ProcessContent(content string, option ContentProcessOption) string {
	if option.RemoveLinks {
		content = linkRegex.ReplaceAllString(content, "")
	}
	if option.RemoveImage {
		content = imgRegex.ReplaceAllString(content, "")
	}
	if option.ConvertToMd {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err == nil {
			doc.Find("script, style, svg").Remove()

			// Remove tags with data URI base64 images
			doc.Find("img").Each(func(i int, s *goquery.Selection) {
				src, exists := s.Attr("src")
				if exists && strings.HasPrefix(strings.ToLower(src), "data:image") {
					s.Remove()
				}
			})

			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				if exists && strings.HasPrefix(strings.ToLower(href), "data:image") {
					s.Remove()
				}
			})

			doc.Find("*").Each(func(i int, s *goquery.Selection) {
				for _, node := range s.Nodes {
					var attrsToRemove []string
					for _, attr := range node.Attr {
						if strings.HasPrefix(attr.Key, "aria-") || strings.HasPrefix(attr.Key, "data-") {
							attrsToRemove = append(attrsToRemove, attr.Key)
						}
					}
					for _, key := range attrsToRemove {
						s.RemoveAttr(key)
					}
				}
			})

			content, _ = doc.Html()
		}

		markdown, err := htmltomarkdown.ConvertString(content)
		if err != nil {
			logrus.Errorf("Error converting HTML to Markdown: %v", err)
			// fallback to original content
			return content
		}
		return markdown
	}
	return content
}
