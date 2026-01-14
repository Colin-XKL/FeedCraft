package parser

import (
	"FeedCraft/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonParser_Parse(t *testing.T) {
	jsonContent := `{
	  "success": true,
	  "data": {
	    "dataList": [
	      {
	        "originalTitle": "Title 1",
	        "url": "http://example.com/1",
	        "publishDateTimeStr": "2023-01-01T00:00:00Z",
	        "summary": "Description 1"
	      },
	      {
	        "originalTitle": "Title 2",
	        "url": "http://example.com/2",
	        "publishDateTimeStr": "2023-01-02T00:00:00Z",
	        "summary": "Description 2"
	      }
	    ]
	  }
	}`

	cfg := &config.JsonParserConfig{
		ItemsIterator: ".data.dataList[]",
		Title:         ".originalTitle",
		Link:          ".url",
		Date:          ".publishDateTimeStr",
		Description:   ".summary",
	}

	parser := &JsonParser{Config: cfg}
	feed, err := parser.Parse([]byte(jsonContent))

	if err != nil {
		t.Logf("Error: %v", err)
	}

	assert.NoError(t, err)
	if feed != nil {
		assert.Len(t, feed.Items, 2)
		if len(feed.Items) > 0 {
			assert.Equal(t, "Title 1", feed.Items[0].Title)
			assert.Equal(t, "http://example.com/1", feed.Items[0].Link)
			assert.Equal(t, "Description 1", feed.Items[0].Description)
		}
	}
}
