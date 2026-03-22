package parser

import (
	"FeedCraft/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/itchyny/gojq"
)

type JsonParsedFields struct {
	Title       string
	Link        string
	Date        string
	Description string
}

type jsonFieldDefinition struct {
	name     string
	selector string
	tmpl     string
	query    *gojq.Query
	parsed   *template.Template
}

type jsonTemplateContext struct {
	Item   interface{}
	Fields JsonParsedFields
}

func ParseJSONItems(rawData interface{}, cfg *config.JsonParserConfig) ([]JsonParsedFields, error) {
	if cfg == nil {
		return nil, fmt.Errorf("parser config is nil")
	}

	itemsArray, err := extractJSONItems(rawData, cfg.ItemsIterator)
	if err != nil {
		return nil, err
	}

	fields, err := compileJSONFieldDefinitions(cfg)
	if err != nil {
		return nil, err
	}

	result := make([]JsonParsedFields, 0, len(itemsArray))
	for _, itemNode := range itemsArray {
		parsedItem, err := parseJSONItemFields(itemNode, fields)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedItem)
	}

	return result, nil
}

func extractJSONItems(rawData interface{}, itemsIterator string) ([]interface{}, error) {
	if itemsIterator == "" || itemsIterator == "." {
		if arr, ok := rawData.([]interface{}); ok {
			return arr, nil
		}
		return []interface{}{rawData}, nil
	}

	query, err := gojq.Parse(itemsIterator)
	if err != nil {
		return nil, fmt.Errorf("failed to parse items_iterator '%s': %w", itemsIterator, err)
	}

	iter := query.Run(rawData)
	var itemsArray []interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, fmt.Errorf("jq execution failed: %w", err)
		}
		if arr, ok := v.([]interface{}); ok {
			itemsArray = append(itemsArray, arr...)
			continue
		}
		itemsArray = append(itemsArray, v)
	}

	return itemsArray, nil
}

func compileJSONFieldDefinitions(cfg *config.JsonParserConfig) ([]jsonFieldDefinition, error) {
	fields := []jsonFieldDefinition{
		{name: "title", selector: cfg.Title, tmpl: cfg.TitleTemplate},
		{name: "link", selector: cfg.Link, tmpl: cfg.LinkTemplate},
		{name: "date", selector: cfg.Date, tmpl: cfg.DateTemplate},
		{name: "description", selector: cfg.Description, tmpl: cfg.DescriptionTemplate},
	}

	for i := range fields {
		if fields[i].selector != "" {
			query, err := gojq.Parse(fields[i].selector)
			if err != nil {
				return nil, fmt.Errorf("invalid %s selector: %w", fields[i].name, err)
			}
			fields[i].query = query
		}

		if fields[i].tmpl != "" {
			tmpl, err := template.New(fields[i].name).
				Option("missingkey=zero").
				Funcs(template.FuncMap{
					"trim": func(value interface{}, cutset string) string {
						return strings.Trim(stringifyTemplateValue(value), cutset)
					},
					"trimSpace": func(value interface{}) string {
						return strings.TrimSpace(stringifyTemplateValue(value))
					},
					"default": func(value interface{}, fallback string) string {
						current := stringifyTemplateValue(value)
						if current == "" {
							return fallback
						}
						return current
					},
				}).
				Parse(fields[i].tmpl)
			if err != nil {
				return nil, fmt.Errorf("invalid %s template: %w", fields[i].name, err)
			}
			fields[i].parsed = tmpl
		}
	}

	return fields, nil
}

func parseJSONItemFields(itemNode interface{}, fields []jsonFieldDefinition) (JsonParsedFields, error) {
	parsed := JsonParsedFields{}
	rawValues := map[string]string{}

	for _, field := range fields {
		value, err := runJSONFieldQuery(field, itemNode)
		if err != nil {
			return JsonParsedFields{}, err
		}
		rawValues[field.name] = value
		switch field.name {
		case "title":
			parsed.Title = value
		case "link":
			parsed.Link = value
		case "date":
			parsed.Date = value
		case "description":
			parsed.Description = value
		}
	}

	ctx := jsonTemplateContext{
		Item: itemNode,
		Fields: JsonParsedFields{
			Title:       rawValues["title"],
			Link:        rawValues["link"],
			Date:        rawValues["date"],
			Description: rawValues["description"],
		},
	}

	for _, field := range fields {
		if field.parsed == nil {
			continue
		}

		rendered, err := renderJSONFieldTemplate(field, ctx)
		if err != nil {
			return JsonParsedFields{}, err
		}

		switch field.name {
		case "title":
			parsed.Title = rendered
		case "link":
			parsed.Link = rendered
		case "date":
			parsed.Date = rendered
		case "description":
			parsed.Description = rendered
		}
	}

	return parsed, nil
}

func runJSONFieldQuery(field jsonFieldDefinition, itemNode interface{}) (string, error) {
	if field.query == nil {
		return "", nil
	}

	iter := field.query.Run(itemNode)
	v, ok := iter.Next()
	if !ok {
		return "", nil
	}
	if err, ok := v.(error); ok {
		return "", fmt.Errorf("failed to extract %s: %w", field.name, err)
	}

	return stringifyTemplateValue(v), nil
}

func renderJSONFieldTemplate(field jsonFieldDefinition, ctx jsonTemplateContext) (string, error) {
	var rendered bytes.Buffer
	if err := field.parsed.Execute(&rendered, ctx); err != nil {
		return "", fmt.Errorf("failed to render %s template: %w", field.name, err)
	}
	return rendered.String(), nil
}

func stringifyTemplateValue(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, err := json.Marshal(v)
	if err == nil {
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}
