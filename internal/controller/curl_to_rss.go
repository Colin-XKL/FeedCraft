package controller

import (
	"FeedCraft/internal/util"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/itchyny/gojq"
)

type JsonFetchReq struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type JsonParseReq struct {
	JsonContent     string `json:"json_content"`
	ListSelector    string `json:"list_selector"`
	TitleSelector   string `json:"title_selector"`
	LinkSelector    string `json:"link_selector"`
	DateSelector    string `json:"date_selector"`
	ContentSelector string `json:"content_selector"`
}

func CurlParseCmd(c *gin.Context) {
	type CurlReq struct {
		CurlCommand string `json:"curl_command"`
	}
	var req CurlReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	parsed, err := util.ParseCurlCommand(req.CurlCommand)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	result := JsonFetchReq{
		Method:  parsed.Method,
		URL:     parsed.URL,
		Headers: parsed.Headers,
		Body:    parsed.Body,
	}

	c.JSON(http.StatusOK, util.APIResponse[JsonFetchReq]{
		StatusCode: 0,
		Data:       result,
	})
}

func CurlFetch(c *gin.Context) {
	var req JsonFetchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	client := resty.New()
	// Disable timeout or set a reasonable one? Using default from rss_tool
	client.SetTimeout(30 * 1000 * 1000 * 1000) // 30s

	r := client.R()
	for k, v := range req.Headers {
		r.SetHeader(k, v)
	}
	if req.Body != "" {
		r.SetBody(req.Body)
	}

	var resp *resty.Response
	var err error

	// Normalize method
	method := strings.ToUpper(req.Method)
	if method == "" {
		method = "GET"
	}

	// Validate URL
	if req.URL == "" {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "URL is required"})
		return
	}
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "invalid URL scheme: must be http or https"})
		return
	}

	// Validate Method
	allowedMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true,
		"PATCH": true, "HEAD": true, "OPTIONS": true,
	}
	if !allowedMethods[method] {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "unsupported method: " + method})
		return
	}

	switch method {
	case "GET":
		resp, err = r.Get(req.URL)
	case "POST":
		resp, err = r.Post(req.URL)
	default:
		resp, err = r.Execute(method, req.URL)
	}

	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "Fetch failed: " + err.Error()})
		return
	}

	// Check for non-2xx status codes
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		body := resp.String()
		if len(body) > 200 {
			body = body[:200] + "..."
		}
		c.JSON(http.StatusOK, util.APIResponse[any]{
			StatusCode: -1,
			Msg:        "Upstream error: " + resp.Status() + " - " + body,
		})
		return
	}

	// Try to format JSON if valid
	var prettyJSON bytes.Buffer
	if json.Unmarshal(resp.Body(), &struct{}{}) == nil { // Simple check if it's JSON object
		if err := json.Indent(&prettyJSON, resp.Body(), "", "  "); err == nil {
			c.JSON(http.StatusOK, util.APIResponse[string]{
				StatusCode: 0,
				Data:       prettyJSON.String(),
			})
			return
		}
	} else {
		// Maybe it's a list or valid json value
		var v interface{}
		if json.Unmarshal(resp.Body(), &v) == nil {
			if err := json.Indent(&prettyJSON, resp.Body(), "", "  "); err == nil {
				c.JSON(http.StatusOK, util.APIResponse[string]{
					StatusCode: 0,
					Data:       prettyJSON.String(),
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, util.APIResponse[string]{
		StatusCode: 0,
		Data:       resp.String(),
	})
}

func CurlParse(c *gin.Context) {
	var req JsonParseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	var input interface{}
	if err := json.Unmarshal([]byte(req.JsonContent), &input); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: "Invalid JSON content: " + err.Error()})
		return
	}

	// Execute List Selector
	listQuery, err := gojq.Parse(req.ListSelector)
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "List selector error: " + err.Error()})
		return
	}

	iter := listQuery.Run(input)
	var rawItems []interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "List extraction error: " + err.Error()})
			return
		}
		// If v is a slice, we iterate it. If it's a single object, we take it.
		// Usually list selector should return an array.
		// If the selector returns multiple results (like .items[]), append them.
		// If it returns one array (like .items), we should iterate that array.

		// gojq behavior:
		// .items returns [obj, obj] -> v is []interface{}
		// .items[] returns obj, obj... -> v is interface{} (called multiple times)

		if arr, ok := v.([]interface{}); ok {
			rawItems = append(rawItems, arr...)
		} else {
			rawItems = append(rawItems, v)
		}
	}

	var parsedItems []ParsedItem

	// Parse other selectors for each item
	// We need to pre-compile selectors for performance, though not critical here
	compile := func(q string) (*gojq.Query, error) {
		if q == "" {
			return nil, nil
		}
		return gojq.Parse(q)
	}

	titleQ, err := compile(req.TitleSelector)
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "invalid TitleSelector: " + err.Error()})
		return
	}
	linkQ, err := compile(req.LinkSelector)
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "invalid LinkSelector: " + err.Error()})
		return
	}
	dateQ, err := compile(req.DateSelector)
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "invalid DateSelector: " + err.Error()})
		return
	}
	contentQ, err := compile(req.ContentSelector)
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "invalid ContentSelector: " + err.Error()})
		return
	}

	runQuery := func(q *gojq.Query, obj interface{}) string {
		if q == nil {
			return ""
		}
		iter := q.Run(obj)
		v, ok := iter.Next()
		if !ok {
			return ""
		}
		if err, ok := v.(error); ok {
			return "Error: " + err.Error()
		}
		// Convert v to string
		if s, ok := v.(string); ok {
			return s
		}
		// If not string, maybe json representation
		b, _ := json.Marshal(v)
		return string(b)
	}

	for _, rawItem := range rawItems {
		item := ParsedItem{}
		item.Title = runQuery(titleQ, rawItem)
		item.Link = runQuery(linkQ, rawItem)
		item.Date = runQuery(dateQ, rawItem)
		item.Content = runQuery(contentQ, rawItem)

		parsedItems = append(parsedItems, item)
	}

	c.JSON(http.StatusOK, util.APIResponse[[]ParsedItem]{
		StatusCode: 0,
		Data:       parsedItems,
	})
}
