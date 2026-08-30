package util

import "strings"

const (
	defaultLLMPromptMaxChars    = 12000
	llmPromptTruncationMark     = "\n\n[...truncated...]\n\n"
	llmPromptTruncationShort    = "[...]"
	llmPromptTruncationEllipsis = "…"
)

// LLMPromptMaxChars returns the max rune count allowed in LLM article payloads.
// Override with FC_LLM_PROMPT_MAX_CHARS. Zero or negative env values fall back to default.
func LLMPromptMaxChars() int {
	envClient := GetEnvClient()
	if envClient != nil {
		if n := envClient.GetInt("LLM_PROMPT_MAX_CHARS"); n > 0 {
			return n
		}
	}
	return defaultLLMPromptMaxChars
}

// TruncateHeadTail keeps the start and end of s when it exceeds maxChars runes.
// Oversized input always includes a truncation marker so the LLM can tell content was cut.
func TruncateHeadTail(s string, maxChars int) string {
	if maxChars <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= maxChars {
		return s
	}

	mark := truncationMarkerForLimit(maxChars)
	remain := maxChars - len(mark)
	if remain < 2 {
		if remain < 0 {
			return string(mark[:maxChars])
		}
		return string(runes[:remain]) + string(mark)
	}

	head := remain * 3 / 5
	if head < 1 {
		head = 1
	}
	tail := remain - head
	if tail < 1 {
		tail = 1
		head = remain - tail
	}

	var b strings.Builder
	b.Grow(maxChars)
	b.WriteString(string(runes[:head]))
	b.WriteString(string(mark))
	b.WriteString(string(runes[len(runes)-tail:]))
	return b.String()
}

func truncationMarkerForLimit(maxChars int) []rune {
	candidates := []string{
		llmPromptTruncationMark,
		llmPromptTruncationShort,
		llmPromptTruncationEllipsis,
	}
	for _, candidate := range candidates {
		mark := []rune(candidate)
		if maxChars >= len(mark)+2 {
			return mark
		}
	}
	mark := []rune(llmPromptTruncationEllipsis)
	if len(mark) > maxChars {
		return mark[:maxChars]
	}
	return mark
}
