package util

import (
	"FeedCraft/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBrowserlessContentUsesFunctionEndpointForNavigationActions(t *testing.T) {
	var captured BrowserFunctionReq
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/function", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&captured))
		_, _ = w.Write([]byte("<html><body>after clicks</body></html>"))
	}))
	defer server.Close()

	t.Setenv("FC_BROWSER_PROVIDER", "browserless-restful")
	t.Setenv("FC_BROWSER_ENDPOINT", server.URL)
	t.Setenv("FC_PUPPETEER_HTTP_ENDPOINT", "")

	content, err := GetBrowserlessContent("https://example.com/list", BrowserlessOptions{
		Timeout:   time.Second,
		WaitUntil: "networkidle2",
		NavigationActions: []config.BrowserNavigationAction{
			{Type: config.BrowserNavigationActionClick, Selector: "#category-a"},
			{Type: config.BrowserNavigationActionClick, Selector: "#category-b"},
			{Type: config.BrowserNavigationActionWait, DurationMs: 500},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "<html><body>after clicks</body></html>", content)
	assert.Equal(t, "https://example.com/list", captured.Context.URL)
	assert.Equal(t, "networkidle2", captured.Context.WaitUntil)
	require.Len(t, captured.Context.NavigationActions, 3)
	assert.Equal(t, "#category-a", captured.Context.NavigationActions[0].Selector)
	assert.Contains(t, captured.Code, "page.click")
}

func TestValidateBrowserNavigationActionsRejectsInvalidAction(t *testing.T) {
	err := ValidateBrowserNavigationActions([]config.BrowserNavigationAction{
		{Type: config.BrowserNavigationActionClick},
	})

	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "selector is required"), err.Error())
}
