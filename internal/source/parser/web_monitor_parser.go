package parser

import (
	"FeedCraft/internal/adapter"
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

// DefaultWebMonitorAIJudgeField is the variable name used to store the AI verdict
// when the user does not specify a custom output field.
const DefaultWebMonitorAIJudgeField = "ai_verdict"

// webMonitorJudgeCaller performs the underlying LLM call for AI judgement.
// It is a package-level variable so tests can stub it out without hitting a
// real LLM or Redis cache. The default implementation reuses
// adapter.CallLLMUsingContext, which caches results by prompt+context hash and
// therefore keeps verdicts (and thus the RSS GUID) stable for identical inputs.
var webMonitorJudgeCaller = adapter.CallLLMUsingContext

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

	doc, values, templates, err := p.prepare(data, p.PageURL)
	if err != nil {
		return nil, err
	}

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
	parser := &WebMonitorParser{Config: cfg, PageURL: pageURL}
	doc, values, templates, err := parser.prepare(data, pageURL)
	if err != nil {
		return nil, err
	}

	rendered, err := renderWebMonitorPreview(doc, values, cfg.KeyFields, templates)
	if err != nil {
		return nil, err
	}
	return rendered.preview, nil
}

func (p *WebMonitorParser) prepare(data []byte, pageURL string) (*goquery.Document, map[string]string, *compiledWebMonitorTemplates, error) {
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
	values["url"] = pageURL

	// Optional AI judgement: derive an extra verdict variable from the
	// extracted values. It runs before key-field validation so the verdict can
	// itself be used as a key field driving the RSS GUID.
	if err := p.applyAIJudge(values); err != nil {
		return nil, nil, nil, err
	}

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

// applyAIJudge runs the optional LLM judgement step and injects the verdict into
// values under the configured output field name.
func (p *WebMonitorParser) applyAIJudge(values map[string]string) error {
	cfg := p.Config.AIJudge
	if cfg == nil || !cfg.Enabled {
		return nil
	}

	prompt := strings.TrimSpace(cfg.Prompt)
	if prompt == "" {
		return fmt.Errorf("ai_judge prompt is required when ai judge is enabled")
	}

	outputField := strings.TrimSpace(cfg.OutputField)
	if outputField == "" {
		outputField = DefaultWebMonitorAIJudgeField
	}
	if outputField == "url" {
		return fmt.Errorf("ai_judge output_field cannot be the reserved variable 'url'")
	}
	if _, exists := p.Config.Extractors[outputField]; exists {
		return fmt.Errorf("ai_judge output_field '%s' conflicts with an existing extractor name", outputField)
	}

	verdict, err := webMonitorJudgeCaller(
		buildWebMonitorJudgePrompt(prompt),
		buildWebMonitorJudgeContext(values),
		util.ContentProcessOption{Temperature: util.LowestLLMTemperaturePtr()},
	)
	if err != nil {
		return fmt.Errorf("ai judge failed: %w", err)
	}

	values[outputField] = normalizeWebMonitorVerdict(verdict)
	return nil
}

func buildWebMonitorJudgePrompt(userPrompt string) string {
	return fmt.Sprintf(`You are a web page change monitoring assistant.

Based on the extracted field values from a web page (provided in the context), follow the user's judgement instruction and output a SHORT verdict label.

User judgement instruction:
%s

Output requirements:
- Output ONLY the verdict label. No explanation, punctuation, quotes, or markdown.
- Keep it short and stable, for example a single word or label such as "available", "unavailable", "yes", or "no".
- Given identical inputs, you MUST always produce the identical verdict.`, userPrompt)
}

// buildWebMonitorJudgeContext renders the extracted values into a deterministic,
// key-sorted text block so the LLM (and its cache) sees stable input.
func buildWebMonitorJudgeContext(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s", key, values[key]))
	}
	return strings.Join(parts, "\n")
}

func normalizeWebMonitorVerdict(raw string) string {
	verdict := strings.TrimSpace(raw)
	if idx := strings.IndexAny(verdict, "\r\n"); idx >= 0 {
		verdict = strings.TrimSpace(verdict[:idx])
	}
	return verdict
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
