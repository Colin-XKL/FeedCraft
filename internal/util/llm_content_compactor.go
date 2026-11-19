package util

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// DefaultLLMCompactorMaxURLLength is the minimum length at which an
	// http/https URL or non-base64 data URI is replaced with a compact
	// placeholder before sending content to an LLM.
	DefaultLLMCompactorMaxURLLength = 128

	llmMediaPlaceholderPrefix = "https://fc.local/m/"
	llmURLPlaceholderPrefix   = "https://fc.local/u/"
)

var safeDataURIMimeTypes = map[string]struct{}{
	"image/png":                {},
	"image/jpeg":               {},
	"image/jpg":                {},
	"image/gif":                {},
	"image/webp":               {},
	"image/avif":               {},
	"image/apng":               {},
	"image/bmp":                {},
	"image/tiff":               {},
	"image/x-icon":             {},
	"image/vnd.microsoft.icon": {},
}

// LLMContentCompactor replaces bulky safe media payloads (such as base64 data
// URIs) and overly long URLs with short deterministic placeholders backed by an
// in-memory dictionary, then restores them verbatim after LLM processing.
//
// Instances are intended to be short-lived and scoped to a single LLM request
// so all stored payloads are eligible for GC immediately after Restore returns.
type LLMContentCompactor struct {
	maxURLLength int
	mediaSeq     int
	urlSeq       int
	placeholders map[string]string
	originals    map[string]string
	reserved     map[string]struct{}
}

// NewLLMContentCompactor creates a compactor using DefaultLLMCompactorMaxURLLength.
func NewLLMContentCompactor() *LLMContentCompactor {
	return NewLLMContentCompactorWithMaxURLLength(DefaultLLMCompactorMaxURLLength)
}

// NewLLMContentCompactorWithMaxURLLength creates a compactor with a custom URL
// length threshold. Non-positive values fall back to DefaultLLMCompactorMaxURLLength.
func NewLLMContentCompactorWithMaxURLLength(maxURLLength int) *LLMContentCompactor {
	if maxURLLength <= 0 {
		maxURLLength = DefaultLLMCompactorMaxURLLength
	}
	return &LLMContentCompactor{
		maxURLLength: maxURLLength,
		placeholders: make(map[string]string),
		originals:    make(map[string]string),
		reserved:     make(map[string]struct{}),
	}
}

// Reserve scans external texts (such as prompt templates) once and marks any
// existing placeholder-shaped strings as unavailable so generated placeholders
// cannot collide with prompt-derived text during Restore.
func (c *LLMContentCompactor) Reserve(texts ...string) {
	if c == nil {
		return
	}
	for _, text := range texts {
		collectExistingPlaceholders(text, c.reserved)
	}
}

// Count returns the number of unique values currently stored in the dictionary.
func (c *LLMContentCompactor) Count() int {
	if c == nil {
		return 0
	}
	return len(c.placeholders)
}

// Compact scans input in linear time, stores matched safe base64/data URIs and
// overly long HTTP(S) URLs in memory, and returns the compacted string.
func (c *LLMContentCompactor) Compact(input string) string {
	if c == nil || input == "" {
		return input
	}

	collectExistingPlaceholders(input, c.reserved)

	var out strings.Builder
	lastFlush := 0
	n := len(input)

	for i := 0; i < n; {
		b := input[i]
		if (b == 'd' || b == 'D') && isURIBoundaryBefore(input, i) && hasPrefixFoldAt(input, i, "data:") {
			end := scanDataURIEnd(input, i)
			candidate := input[i:end]
			if shouldCompactDataURI(candidate, c.maxURLLength) {
				out.WriteString(input[lastFlush:i])
				out.WriteString(c.register(candidate, llmMediaPlaceholderPrefix))
				lastFlush = end
				i = end
				continue
			}
			i = end
			continue
		}

		if (b == 'h' || b == 'H') && isURIBoundaryBefore(input, i) && hasHTTPPrefixAt(input, i) {
			end := scanHTTPURIEnd(input, i)
			candidate := input[i:end]
			if shouldCompactHTTPURL(candidate, c.maxURLLength) {
				out.WriteString(input[lastFlush:i])
				out.WriteString(c.register(candidate, llmURLPlaceholderPrefix))
				lastFlush = end
				i = end
				continue
			}
			i = end
			continue
		}

		i++
	}

	if lastFlush == 0 {
		return input
	}
	out.WriteString(input[lastFlush:])
	return out.String()
}

// Restore replaces all placeholders present in output with their original values.
func (c *LLMContentCompactor) Restore(output string) string {
	if c == nil || output == "" || len(c.placeholders) == 0 {
		return output
	}

	if len(c.placeholders) == 1 {
		for placeholder, original := range c.placeholders {
			return strings.ReplaceAll(output, placeholder, original)
		}
	}

	keys := make([]string, 0, len(c.placeholders))
	for k := range c.placeholders {
		keys = append(keys, k)
	}
	// Sort by descending key length so longer placeholders (e.g. .../10) are
	// matched before shorter prefixes (e.g. .../1) by strings.NewReplacer.
	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) > len(keys[j])
		}
		return keys[i] < keys[j]
	})

	oldnew := make([]string, 0, len(keys)*2)
	for _, k := range keys {
		oldnew = append(oldnew, k, c.placeholders[k])
	}
	return strings.NewReplacer(oldnew...).Replace(output)
}

func (c *LLMContentCompactor) register(original string, prefix string) string {
	if existing, ok := c.originals[original]; ok {
		return existing
	}

	var placeholder string
	for {
		var seq int
		if prefix == llmMediaPlaceholderPrefix {
			c.mediaSeq++
			seq = c.mediaSeq
		} else {
			c.urlSeq++
			seq = c.urlSeq
		}
		candidate := fmt.Sprintf("%s%d", prefix, seq)
		if _, used := c.placeholders[candidate]; used {
			continue
		}
		if _, collision := c.reserved[candidate]; collision {
			continue
		}
		placeholder = candidate
		break
	}

	c.placeholders[placeholder] = original
	c.originals[original] = placeholder
	return placeholder
}

// collectExistingPlaceholders scans text once per placeholder prefix and records
// every numeric placeholder prefix that occurs in text (including prefixes of
// longer digit runs, matching strings.Contains semantics).
func collectExistingPlaceholders(text string, reserved map[string]struct{}) {
	if text == "" || reserved == nil {
		return
	}
	collectPrefixPlaceholders(text, llmMediaPlaceholderPrefix, reserved)
	collectPrefixPlaceholders(text, llmURLPlaceholderPrefix, reserved)
}

func collectPrefixPlaceholders(text string, prefix string, reserved map[string]struct{}) {
	searchFrom := 0
	for searchFrom < len(text) {
		rel := strings.Index(text[searchFrom:], prefix)
		if rel < 0 {
			return
		}
		digitStart := searchFrom + rel + len(prefix)
		i := digitStart
		for i < len(text) && text[i] >= '0' && text[i] <= '9' {
			i++
			reserved[text[searchFrom+rel:i]] = struct{}{}
		}
		searchFrom = digitStart
	}
}

func isURIBoundaryBefore(s string, idx int) bool {
	if idx <= 0 {
		return true
	}
	prev := s[idx-1]
	isAlphaNum := (prev >= 'a' && prev <= 'z') ||
		(prev >= 'A' && prev <= 'Z') ||
		(prev >= '0' && prev <= '9') ||
		prev == '_'
	return !isAlphaNum
}

func hasPrefixFoldAt(s string, idx int, prefix string) bool {
	if idx+len(prefix) > len(s) {
		return false
	}
	return strings.EqualFold(s[idx:idx+len(prefix)], prefix)
}

func hasHTTPPrefixAt(s string, idx int) bool {
	return hasPrefixFoldAt(s, idx, "https://") ||
		hasPrefixFoldAt(s, idx, "http://") ||
		hasPrefixFoldAt(s, idx, `https:\/\/`) ||
		hasPrefixFoldAt(s, idx, `http:\/\/`)
}

func scanDataURIEnd(s string, start int) int {
	n := len(s)
	quote := byte(0)
	if start > 0 && (s[start-1] == '"' || s[start-1] == '\'') {
		quote = s[start-1]
	}

	i := start
	parenDepth := 0
	for i < n {
		b := s[i]
		if b == '\\' {
			if i+1 < n && s[i+1] == '/' {
				i += 2
				continue
			}
			break
		}
		if b <= 0x20 || b == 0x7f || b == '<' || b == '>' || b == '`' {
			break
		}
		if quote != 0 {
			if b == quote {
				break
			}
			i++
			continue
		}
		if b == '"' || b == '\'' {
			break
		}
		if b == '(' {
			parenDepth++
		} else if b == ')' {
			if parenDepth == 0 {
				break
			}
			parenDepth--
		}
		i++
	}

	// Handle quoted base64 data URIs that contain wrapped line breaks or spaces
	// inside the base64 payload (e.g. HTML generated by email clients).
	if quote != 0 && i < n && isASCIIWhitespace(s[i]) && isBase64DataURIHeader(s[start:i]) {
		j := i
		for j < n {
			b := s[j]
			if b == quote {
				return j
			}
			if b == '\\' && j+1 < n && s[j+1] == '/' {
				j += 2
				continue
			}
			if !isBase64PayloadByte(b) && !isASCIIWhitespace(b) {
				break
			}
			j++
		}
	}

	return i
}

func scanHTTPURIEnd(s string, start int) int {
	n := len(s)
	quote := byte(0)
	if start > 0 && (s[start-1] == '"' || s[start-1] == '\'') {
		quote = s[start-1]
	}

	i := start
	parenDepth := 0
	bracketDepth := 0
	for i < n {
		b := s[i]
		if b == '\\' {
			if i+1 < n && s[i+1] == '/' {
				i += 2
				continue
			}
			break
		}
		if b <= 0x20 || b == 0x7f || b == '<' || b == '>' || b == '`' {
			break
		}
		if quote != 0 {
			if b == quote {
				break
			}
			i++
			continue
		}
		if b == '"' || b == '\'' {
			break
		}
		switch b {
		case '(':
			parenDepth++
		case ')':
			if parenDepth == 0 {
				goto done
			}
			parenDepth--
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth == 0 {
				goto done
			}
			bracketDepth--
		}
		i++
	}

done:
	if quote == 0 {
		for i > start {
			last := s[i-1]
			if last == '.' || last == ',' || last == ';' {
				i--
				continue
			}
			break
		}
	}
	return i
}

func isASCIIWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f' || b == '\v'
}

func isBase64PayloadByte(b byte) bool {
	return (b >= 'A' && b <= 'Z') ||
		(b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') ||
		b == '+' || b == '/' || b == '=' || b == '-' || b == '_' || b == '%'
}

func isBase64DataURIHeader(s string) bool {
	comma := strings.IndexByte(s, ',')
	if comma < 0 {
		return false
	}
	header := strings.ToLower(s[:comma])
	return strings.Contains(header, ";base64")
}

func isSafeMediaDataURIMime(mime string) bool {
	if _, ok := safeDataURIMimeTypes[mime]; ok {
		return true
	}
	return strings.HasPrefix(mime, "audio/") ||
		strings.HasPrefix(mime, "video/") ||
		strings.HasPrefix(mime, "font/")
}

func shouldCompactDataURI(candidate string, maxURLLength int) bool {
	comma := strings.IndexByte(candidate, ',')
	if comma < 0 || comma == len(candidate)-1 {
		return false
	}
	header := strings.ToLower(strings.ReplaceAll(candidate[len("data:"):comma], `\/`, `/`))
	mime := header
	if semi := strings.IndexByte(header, ';'); semi >= 0 {
		mime = header[:semi]
	}
	mime = strings.TrimSpace(mime)
	if !isSafeMediaDataURIMime(mime) {
		return false
	}
	if strings.Contains(header, ";base64") {
		return strings.TrimSpace(candidate[comma+1:]) != ""
	}
	return len(candidate) >= maxURLLength
}

func shouldCompactHTTPURL(candidate string, maxURLLength int) bool {
	if len(candidate) < maxURLLength {
		return false
	}
	lower := strings.ToLower(candidate)
	if strings.HasPrefix(lower, llmMediaPlaceholderPrefix) || strings.HasPrefix(lower, llmURLPlaceholderPrefix) {
		return false
	}
	return true
}
