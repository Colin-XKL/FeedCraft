package util

import (
	"regexp"
)

var (
	linkRegex = regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"[^>]*>.*?</a>`)
	imgRegex  = regexp.MustCompile(`<img\s+(?:[^>]*?\s+)?src="([^"]*)"[^>]*>`)
)

type ContentProcessOption struct {
	RemoveLinks bool
	RemoveImage bool
	ConvertToMd bool
	Temperature *float64
}

const LowestLLMTemperature = 0.0

func LowestLLMTemperaturePtr() *float64 {
	temperature := LowestLLMTemperature
	return &temperature
}

func ProcessContent(content string, option ContentProcessOption) string {
	if option.RemoveLinks {
		content = linkRegex.ReplaceAllString(content, "")
	}
	if option.RemoveImage {
		content = imgRegex.ReplaceAllString(content, "")
	}
	if option.ConvertToMd {
		if md := HTMLToMarkdown(content, ""); md != "" {
			return md
		}
		return content
	}
	return content
}
