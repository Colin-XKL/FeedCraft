package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidateURLAllowsPrivateAndLoopbackAddresses(t *testing.T) {
	t.Parallel()

	cases := []string{
		"http://127.0.0.1/article/1.html",
		"http://192.168.5.2:10089/article/1.html",
		"http://10.0.0.8/page",
		"http://172.16.0.4/page",
	}
	for _, raw := range cases {
		if err := validateURL(raw); err != nil {
			t.Errorf("validateURL(%q) = %v, want nil (admin HTML fetch should allow local/private sources)", raw, err)
		}
	}
}

func TestValidateURLRejectsLinkLocalMetadataAddresses(t *testing.T) {
	t.Parallel()

	err := validateURL("http://169.254.169.254/latest/meta-data/")
	if err == nil {
		t.Fatal("validateURL allowed link-local metadata address, want rejection")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "forbidden") {
		t.Fatalf("error = %q, want a forbidden message", err)
	}
}

func TestValidateURLRejectsNonHTTPSchemes(t *testing.T) {
	t.Parallel()

	err := validateURL("file:///etc/passwd")
	if err == nil {
		t.Fatal("validateURL allowed file scheme, want rejection")
	}
}

func TestHtmlFetchAllowsLoopbackSource(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><head><title>Local Mock</title></head><body><h1>ok</h1></body></html>`))
	}))
	defer upstream.Close()

	recorder := performHtmlFetchRequest(t, FetchReq{URL: upstream.URL})

	var response struct {
		Code int    `json:"code"`
		Data string `json:"data"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if recorder.Code != http.StatusOK || response.Code != 0 {
		t.Fatalf("expected fetch success for loopback URL, got http %d code %d msg %q body %s",
			recorder.Code, response.Code, response.Msg, recorder.Body.String())
	}
	if !strings.Contains(response.Data, "Local Mock") {
		t.Fatalf("expected fetched HTML, got %q", response.Data)
	}
}

func TestHtmlFetchRejectsLinkLocalWithReadableMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performHtmlFetchRequest(t, FetchReq{URL: "http://169.254.169.254/latest/meta-data/"})

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code == 0 {
		t.Fatalf("expected validation failure, got success body %s", recorder.Body.String())
	}
	msg := strings.ToLower(response.Msg)
	if strings.Contains(msg, "request failed with status code") {
		t.Fatalf("msg should be human-readable, got %q", response.Msg)
	}
	if !strings.Contains(msg, "forbidden") && !strings.Contains(msg, "not allowed") {
		t.Fatalf("msg = %q, want a readable SSRF/forbidden explanation", response.Msg)
	}
}

func performHtmlFetchRequest(t *testing.T, body FetchReq) *httptest.ResponseRecorder {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/fetch", HtmlFetch)
	req, err := http.NewRequest(http.MethodPost, "/fetch", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	return recorder
}
