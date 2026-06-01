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

// paragraphTagRE matches <p> elements so we can normalize line breaks inside them.
var paragraphTagRE = regexp.MustCompile(`(?is)<p(\s[^>]*)?>(.*?)</p>`)

func Markdown2HTML(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	renderResult := markdown.Render(doc, renderer)
	return insertLineBreaksInParagraphHTML(string(renderResult))
}

func Html2Markdown(text string, domain *string) string {
	text = insertLineBreaksInParagraphHTML(text)

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			table.NewTablePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)
	var convertOptions []converter.ConvertOptionFunc
	if domain != nil {
		convertOptions = append(convertOptions, converter.WithDomain(*domain))
	}

	mdStr, err := conv.ConvertString(text, convertOptions...)

	if err != nil {
		logrus.Errorf("convert html to markdown err: %v", err)
	}
	return mdStr
}

// insertLineBreaksInParagraphHTML converts literal newlines inside <p> tags into <br>,
// so soft line breaks survive browser rendering and HTML↔Markdown round-trips.
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
	lower := strings.ToLower(before)
	return !strings.HasSuffix(lower, "<br>") && !strings.HasSuffix(lower, "<br/>") && !strings.HasSuffix(lower, "<br />")
}
