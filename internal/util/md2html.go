package util

import (
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/sirupsen/logrus"
)

func Markdown2HTML(md string) string {
	// create Markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	renderResult := markdown.Render(doc, renderer)
	return string(renderResult)
}

func Html2Markdown(text string, domain *string) string {
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
