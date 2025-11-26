package util

import (
	"regexp"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
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
