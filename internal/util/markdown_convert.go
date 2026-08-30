package util

import (
	"regexp"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/sirupsen/logrus"
)

var (
	paragraphTagRE          = regexp.MustCompile(`(?is)<p(\s[^>]*)?>(.*?)</p>`)
	consecutiveBlankLinesRE = regexp.MustCompile(`(?:\r?\n[ \t]*){3,}`)
	brTagSuffixRE           = regexp.MustCompile(`(?i)<br\b[^>]*>\s*$`)
	brTagPrefixRE           = regexp.MustCompile(`(?i)^\s*<br\b[^>]*>`)

	blockStartRE      = regexp.MustCompile("^(?:#{1,6}\\s|(?:\\*|-|\\+)\\s|\\d+\\.\\s|>|```|\\|)")
	orderedListItemRE = regexp.MustCompile(`^\d+\.\s`)
)

func MarkdownToHTML(md string) string {
	md = normalizeMarkdownForRender(md)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func HTMLToMarkdown(htmlContent string, domain string) string {
	htmlContent = insertLineBreaksInParagraphHTML(htmlContent)

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			table.NewTablePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)
	var convertOptions []converter.ConvertOptionFunc
	if domain != "" {
		convertOptions = append(convertOptions, converter.WithDomain(domain))
	}

	mdStr, err := conv.ConvertString(htmlContent, convertOptions...)
	if err != nil {
		logrus.Errorf("convert html to markdown err: %v", err)
	}
	return collapseConsecutiveBlankLines(mdStr)
}

func normalizeMarkdownForRender(md string) string {
	return processOutsideFencedCodeBlocks(md, func(segment string) string {
		segment = collapseConsecutiveBlankLines(segment)
		return applySingleNewlineHardBreaks(segment)
	})
}

func collapseConsecutiveBlankLines(md string) string {
	return consecutiveBlankLinesRE.ReplaceAllString(md, "\n\n")
}

func applySingleNewlineHardBreaks(md string) string {
	if md == "" {
		return md
	}

	lines := strings.Split(md, "\n")
	for i := 0; i < len(lines)-1; i++ {
		current := lines[i]
		next := lines[i+1]
		if current == "" || next == "" {
			continue
		}
		if isListContinuation(current, next) {
			continue
		}
		if blockStartRE.MatchString(next) && !isSameBlockContinuation(current, next) {
			continue
		}
		lines[i] = strings.TrimRight(current, " ") + "  "
	}
	return strings.Join(lines, "\n")
}

func isSameBlockContinuation(prev, next string) bool {
	prevTrim := strings.TrimSpace(prev)
	nextTrim := strings.TrimSpace(next)
	return strings.HasPrefix(prevTrim, "> ") && strings.HasPrefix(nextTrim, "> ")
}

func isListContinuation(prev, next string) bool {
	if strings.TrimSpace(next) == "" {
		return false
	}
	indent := len(next) - len(strings.TrimLeft(next, " \t"))
	if indent < 2 {
		return false
	}
	prevTrim := strings.TrimSpace(prev)
	if strings.HasPrefix(prevTrim, "- ") || strings.HasPrefix(prevTrim, "* ") || strings.HasPrefix(prevTrim, "+ ") {
		return true
	}
	return orderedListItemRE.MatchString(prevTrim)
}

func processOutsideFencedCodeBlocks(md string, process func(string) string) string {
	var result strings.Builder
	remaining := md

	for {
		open := strings.Index(remaining, "```")
		if open == -1 {
			result.WriteString(process(remaining))
			break
		}

		result.WriteString(process(remaining[:open]))
		rest := remaining[open+3:]
		close := strings.Index(rest, "```")
		if close == -1 {
			result.WriteString(remaining[open:])
			break
		}

		closeEnd := open + 3 + close + 3
		result.WriteString(remaining[open:closeEnd])
		remaining = remaining[closeEnd:]
	}

	return result.String()
}

func insertLineBreaksInParagraphHTML(html string) string {
	return paragraphTagRE.ReplaceAllStringFunc(html, func(match string) string {
		subs := paragraphTagRE.FindStringSubmatch(match)
		attrs := subs[1]
		inner := subs[2]
		return "<p" + attrs + ">" + insertLineBreaksInHTMLFragment(inner) + "</p>"
	})
}

func insertLineBreaksInHTMLFragment(fragment string) string {
	if !strings.Contains(fragment, "\n") {
		return fragment
	}

	var b strings.Builder
	b.Grow(len(fragment) + 16)
	segmentStart := 0

	for i := 0; i < len(fragment); i++ {
		if fragment[i] != '\n' {
			continue
		}
		before := fragment[segmentStart:i]
		after := fragment[i+1:]
		if shouldInsertLineBreak(before, after) {
			b.WriteString(before)
			b.WriteString("<br>\n")
			segmentStart = i + 1
		}
	}

	b.WriteString(fragment[segmentStart:])
	return b.String()
}

func shouldInsertLineBreak(before, after string) bool {
	before = strings.TrimRight(before, " \t\r\n")
	after = strings.TrimLeft(after, " \t\r\n")
	if before == "" || after == "" {
		return false
	}
	return !brTagSuffixRE.MatchString(before) && !brTagPrefixRE.MatchString(after)
}
