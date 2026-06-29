package util

import (
	"regexp"
)

var (
	linkRegex = regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"[^>]*>.*?</a>`)
	imgRegex  = regexp.MustCompile(`<img\s+(?:[^>]*?\s+)?src="([^"]*)"[^>]*>`)

	// base64ImgTagRegex matches <img>/<source> tags whose attributes embed a
	// base64-encoded data URI (e.g. src="data:image/png;base64,....").
	// Quote style and attribute order are irrelevant since we only require the
	// presence of a base64 data URI inside the tag.
	base64ImgTagRegex = regexp.MustCompile(`(?is)<(?:img|source)\b[^>]*?\bdata:[^>]*?;base64,[^>]*?>`)
	// base64MarkdownImgRegex matches Markdown image syntax that embeds a base64
	// data URI, e.g. ![alt](data:image/png;base64,....). This acts as a safety
	// net for content that is already in Markdown form.
	base64MarkdownImgRegex = regexp.MustCompile(`(?is)!\[[^\]]*\]\(\s*data:[^)]*?;base64,[^)]*?\)`)
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

// RemoveBase64Images strips inline base64-encoded images from the given
// content. Such images carry no useful signal for LLM processing but can
// dominate the token budget, so we drop them before sending content to the LLM.
// It handles both HTML (<img>/<source> data URIs) and Markdown image syntax.
func RemoveBase64Images(content string) string {
	if content == "" {
		return content
	}
	content = base64ImgTagRegex.ReplaceAllString(content, "")
	content = base64MarkdownImgRegex.ReplaceAllString(content, "")
	return content
}

func ProcessContent(content string, option ContentProcessOption) string {
	if option.RemoveLinks {
		content = linkRegex.ReplaceAllString(content, "")
	}
	if option.RemoveImage {
		content = imgRegex.ReplaceAllString(content, "")
		// imgRegex only matches double-quoted src attributes, so explicitly
		// drop any remaining base64-encoded images regardless of quote style.
		content = RemoveBase64Images(content)
	}
	if option.ConvertToMd {
		if md := HTMLToMarkdown(content, ""); md != "" {
			return md
		}
		return content
	}
	return content
}
