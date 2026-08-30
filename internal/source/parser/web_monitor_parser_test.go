package parser

import (
	"FeedCraft/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebMonitorParser_Parse_KeyFieldsDriveGUID(t *testing.T) {
	html := `
	<html>
	  <head><title>Product Page</title></head>
	  <body>
	    <div class="price"> $399 </div>
	    <div class="stock">In Stock</div>
	    <h1 class="title">PS5 Console</h1>
	  </body>
	</html>`

	cfg := &config.WebMonitorParserConfig{
		Extractors: map[string]string{
			"price": ".price",
			"stock": ".stock",
			"title": ".title",
		},
		KeyFields:           []string{"price"},
		TitleTemplate:       "价格 {{.price}}",
		DescriptionTemplate: "库存 {{.stock}}",
		ContentTemplate:     "{{.title}} - {{.price}} - {{.stock}} - {{.url}}",
	}

	parser := &WebMonitorParser{Config: cfg, PageURL: "https://example.com/product"}
	feed, err := parser.Parse([]byte(html))

	require.NoError(t, err)
	require.NotNil(t, feed)
	require.Len(t, feed.Articles, 1)
	assert.Equal(t, "Product Page", feed.Title)
	assert.Equal(t, "价格 $399", feed.Articles[0].Title)
	assert.Equal(t, "库存 In Stock", feed.Articles[0].Description)
	assert.Equal(t, "https://example.com/product", feed.Articles[0].Link)
	assert.Equal(t, "PS5 Console - $399 - In Stock - https://example.com/product", feed.Articles[0].Content)

	feedSame, err := parser.Parse([]byte(html))
	require.NoError(t, err)
	assert.Equal(t, feed.Articles[0].Id, feedSame.Articles[0].Id)

	htmlStockChanged := `
	<html>
	  <head><title>Product Page</title></head>
	  <body>
	    <div class="price">$399</div>
	    <div class="stock">Only 2 left</div>
	    <h1 class="title">PS5 Console</h1>
	  </body>
	</html>`
	feedStockChanged, err := parser.Parse([]byte(htmlStockChanged))
	require.NoError(t, err)
	assert.Equal(t, feed.Articles[0].Id, feedStockChanged.Articles[0].Id)

	htmlPriceChanged := `
	<html>
	  <head><title>Product Page</title></head>
	  <body>
	    <div class="price">$499</div>
	    <div class="stock">In Stock</div>
	    <h1 class="title">PS5 Console</h1>
	  </body>
	</html>`
	feedPriceChanged, err := parser.Parse([]byte(htmlPriceChanged))
	require.NoError(t, err)
	assert.NotEqual(t, feed.Articles[0].Id, feedPriceChanged.Articles[0].Id)
}

func TestWebMonitorParser_Parse_URLInjectedIntoTemplateAndGUID(t *testing.T) {
	html := `<html><body><div class="price">$399</div></body></html>`

	cfg := &config.WebMonitorParserConfig{
		Extractors:      map[string]string{"price": ".price"},
		KeyFields:       []string{"price"},
		ContentTemplate: "链接：{{.url}}",
	}

	parserA := &WebMonitorParser{Config: cfg, PageURL: "https://site-a.com/product"}
	parserB := &WebMonitorParser{Config: cfg, PageURL: "https://site-b.com/product"}

	feedA, err := parserA.Parse([]byte(html))
	require.NoError(t, err)
	feedB, err := parserB.Parse([]byte(html))
	require.NoError(t, err)

	// URL should be injected into content template
	assert.Equal(t, "链接：https://site-a.com/product", feedA.Articles[0].Content)
	assert.Equal(t, "链接：https://site-b.com/product", feedB.Articles[0].Content)

	// URL should be the article link
	assert.Equal(t, "https://site-a.com/product", feedA.Articles[0].Link)
	assert.Equal(t, "https://site-b.com/product", feedB.Articles[0].Link)

	// Different URLs must produce different GUIDs even when key field values are identical
	assert.NotEqual(t, feedA.Articles[0].Id, feedB.Articles[0].Id)
}

func TestWebMonitorParser_Parse_KeyFieldOrderDoesNotMatter(t *testing.T) {
	html := `<html><body><div class="price">$399</div><div class="stock">In Stock</div></body></html>`

	parserA := &WebMonitorParser{Config: &config.WebMonitorParserConfig{
		Extractors: map[string]string{"price": ".price", "stock": ".stock"},
		KeyFields:  []string{"price", "stock"},
	}}
	parserB := &WebMonitorParser{Config: &config.WebMonitorParserConfig{
		Extractors: map[string]string{"stock": ".stock", "price": ".price"},
		KeyFields:  []string{"stock", "price"},
	}}

	feedA, err := parserA.Parse([]byte(html))
	require.NoError(t, err)
	feedB, err := parserB.Parse([]byte(html))
	require.NoError(t, err)
	require.Len(t, feedA.Articles, 1)
	require.Len(t, feedB.Articles, 1)
	assert.Equal(t, feedA.Articles[0].Id, feedB.Articles[0].Id)
}

func TestWebMonitorParser_Parse_DefaultContentFallsBackToDescription(t *testing.T) {
	parser := &WebMonitorParser{Config: &config.WebMonitorParserConfig{
		Extractors:          map[string]string{"price": ".price"},
		KeyFields:           []string{"price"},
		DescriptionTemplate: "价格 {{.price}}",
	}}

	feed, err := parser.Parse([]byte(`<html><body><div class="price"> $399 </div></body></html>`))
	require.NoError(t, err)
	require.Len(t, feed.Articles, 1)
	assert.Equal(t, "价格 $399", feed.Articles[0].Description)
	assert.Equal(t, "价格 $399", feed.Articles[0].Content)
}

func TestWebMonitorParser_Parse_DefaultContentFallbackSummary(t *testing.T) {
	parser := &WebMonitorParser{Config: &config.WebMonitorParserConfig{
		Extractors: map[string]string{"price": ".price", "stock": ".stock"},
		KeyFields:  []string{"price"},
	}}

	feed, err := parser.Parse([]byte(`<html><body><div class="price"> $399 </div><div class="stock"> In Stock </div></body></html>`))
	require.NoError(t, err)
	require.Len(t, feed.Articles, 1)
	assert.Contains(t, feed.Articles[0].Content, "price: $399")
	assert.Contains(t, feed.Articles[0].Content, "stock: In Stock")
}

func TestWebMonitorParser_Parse_MissingKeyFieldFails(t *testing.T) {
	parser := &WebMonitorParser{Config: &config.WebMonitorParserConfig{
		Extractors: map[string]string{"price": ".price"},
		KeyFields:  []string{"stock"},
	}}

	feed, err := parser.Parse([]byte(`<html><body><div class="price">$399</div></body></html>`))
	assert.Error(t, err)
	assert.Nil(t, feed)
	assert.Contains(t, err.Error(), "key field 'stock' is not defined in extractors")
}

func TestWebMonitorParser_Parse_InvalidTemplateFails(t *testing.T) {
	parser := &WebMonitorParser{Config: &config.WebMonitorParserConfig{
		Extractors:    map[string]string{"price": ".price"},
		KeyFields:     []string{"price"},
		TitleTemplate: "{{.price}",
	}}

	feed, err := parser.Parse([]byte(`<html><body><div class="price">$399</div></body></html>`))
	assert.Error(t, err)
	assert.Nil(t, feed)
	assert.Contains(t, err.Error(), "invalid title template")
}
