package parser

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type WebMonitorParser struct {
	Config  *config.WebMonitorParserConfig
	PageURL string
}

type webMonitorTemplateContext map[string]string

type compiledWebMonitorTemplates struct {
	title       *template.Template
	description *template.Template
	content     *template.Template
}

func (p *WebMonitorParser) Parse(data []byte) (*model.CraftFeed, error) {
	if p == nil || p.Config == nil {
		return nil, fmt.Errorf("parser config is nil")
	}

	doc, values, templates, err := p.prepare(data)
	if err != nil {
		return nil, err
	}

	values["url"] = p.PageURL

	rendered, err := renderWebMonitorPreview(doc, values, p.Config.KeyFields, templates)
	if err != nil {
		return nil, err
	}

	return rendered.feed, nil
}

type WebMonitorPreviewResult struct {
	Values      map[string]string `json:"values"`
	KeyFields   []string          `json:"key_fields"`
	GUID        string            `json:"guid"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Content     string            `json:"content"`
	FeedTitle   string            `json:"feed_title"`
	FeedLink    string            `json:"feed_link"`
}

type webMonitorRendered struct {
	feed    *model.CraftFeed
	preview *WebMonitorPreviewResult
}

func PreviewWebMonitor(data []byte, cfg *config.WebMonitorParserConfig, pageURL string) (*WebMonitorPreviewResult, error) {
	parser := &WebMonitorParser{Config: cfg}
	doc, values, templates, err := parser.prepare(data)
	if err != nil {
		return nil, err
	}
	values["url"] = pageURL

	rendered, err := renderWebMonitorPreview(doc, values, cfg.KeyFields, templates)
	if err != nil {
		return nil, err
	}
	return rendered.preview, nil
}

func (p *WebMonitorParser) prepare(data []byte) (*goquery.Document, map[string]string, *compiledWebMonitorTemplates, error) {
	if len(p.Config.Extractors) == 0 {
		return nil, nil, nil, fmt.Errorf("extractors are required")
	}

	if len(p.Config.KeyFields) == 0 {
		return nil, nil, nil, fmt.Errorf("key_fields are required")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse html: %w", err)
	}

	values := extractWebMonitorValues(doc, p.Config.Extractors)
	for _, field := range p.Config.KeyFields {
		if _, ok := values[field]; !ok {
			return nil, nil, nil, fmt.Errorf("key field '%s' is not defined in extractors", field)
		}
	}

	templates, err := compileWebMonitorTemplates(p.Config)
	if err != nil {
		return nil, nil, nil, err
	}

	return doc, values, templates, nil
}

func extractWebMonitorValues(doc *goquery.Document, extractors map[string]string) map[string]string {
	values := make(map[string]string, len(extractors))
	for name, selector := range extractors {
		trimmedSelector := strings.TrimSpace(selector)
		if trimmedSelector == "" {
			values[name] = ""
			continue
		}
		values[name] = strings.TrimSpace(doc.Find(trimmedSelector).First().Text())
	}
	return values
}

func compileWebMonitorTemplates(cfg *config.WebMonitorParserConfig) (*compiledWebMonitorTemplates, error) {
	title, err := compileWebMonitorTemplate("title", cfg.TitleTemplate)
	if err != nil {
		return nil, fmt.Errorf("invalid title template: %w", err)
	}

	description, err := compileWebMonitorTemplate("description", cfg.DescriptionTemplate)
	if err != nil {
		return nil, fmt.Errorf("invalid description template: %w", err)
	}

	content, err := compileWebMonitorTemplate("content", cfg.ContentTemplate)
	if err != nil {
		return nil, fmt.Errorf("invalid content template: %w", err)
	}

	return &compiledWebMonitorTemplates{
		title:       title,
		description: description,
		content:     content,
	}, nil
}

func compileWebMonitorTemplate(name, src string) (*template.Template, error) {
	if strings.TrimSpace(src) == "" {
		return nil, nil
	}
	return template.New(name).Option("missingkey=zero").Parse(src)
}

func renderWebMonitorTemplate(tmpl *template.Template, ctx webMonitorTemplateContext) (string, error) {
	if tmpl == nil {
		return "", nil
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, ctx); err != nil {
		return "", err
	}
	return rendered.String(), nil
}

func renderWebMonitorContent(tmpl *template.Template, description string, ctx webMonitorTemplateContext) (string, error) {
	if tmpl != nil {
		rendered, err := renderWebMonitorTemplate(tmpl, ctx)
		if err != nil {
			return "", fmt.Errorf("failed to render content template: %w", err)
		}
		return rendered, nil
	}
	if description != "" {
		return description, nil
	}
	return defaultWebMonitorContent(ctx), nil
}

func renderWebMonitorPreview(doc *goquery.Document, values map[string]string, keyFields []string, templates *compiledWebMonitorTemplates) (*webMonitorRendered, error) {
	ctx := webMonitorTemplateContext(values)
	if _, ok := ctx["url"]; !ok {
		ctx["url"] = ""
	}

	title, err := renderWebMonitorTemplate(templates.title, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to render title template: %w", err)
	}

	description, err := renderWebMonitorTemplate(templates.description, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to render description template: %w", err)
	}

	content, err := renderWebMonitorContent(templates.content, description, ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	pageURL := ctx["url"]
	guid := buildWebMonitorGUID(pageURL, keyFields, values)
	feedTitle := strings.TrimSpace(doc.Find("title").First().Text())
	if feedTitle == "" {
		feedTitle = pageURL
	}
	if feedTitle == "" {
		feedTitle = "Web Monitor"
	}

	article := &model.CraftArticle{
		Title:       title,
		Link:        pageURL,
		Description: description,
		Id:          guid,
		Created:     now,
		Updated:     now,
		Content:     content,
	}

	previewValues := make(map[string]string, len(values))
	for key, value := range values {
		previewValues[key] = value
	}

	return &webMonitorRendered{
		feed: &model.CraftFeed{
			Title:       feedTitle,
			Link:        pageURL,
			Description: description,
			Created:     now,
			Updated:     now,
			Articles:    []*model.CraftArticle{article},
		},
		preview: &WebMonitorPreviewResult{
			Values:      previewValues,
			KeyFields:   append([]string(nil), keyFields...),
			GUID:        guid,
			Title:       title,
			Description: description,
			Content:     content,
			FeedTitle:   feedTitle,
			FeedLink:    pageURL,
		},
	}, nil
}

func defaultWebMonitorContent(ctx webMonitorTemplateContext) string {
	keys := make([]string, 0, len(ctx))
	for key := range ctx {
		if key == "url" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys)+1)
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s", key, ctx[key]))
	}
	if ctx["url"] != "" {
		parts = append(parts, fmt.Sprintf("url: %s", ctx["url"]))
	}
	return strings.Join(parts, "\n")
}

func buildWebMonitorGUID(pageURL string, keyFields []string, values map[string]string) string {
	sortedKeys := append([]string(nil), keyFields...)
	sort.Strings(sortedKeys)

	parts := make([]string, 0, len(sortedKeys)+1)
	parts = append(parts, pageURL)
	for _, key := range sortedKeys {
		parts = append(parts, key+"="+values[key])
	}
	return util.GetPasswordMD5Hash(strings.Join(parts, "\n"))
}
