package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCheckBrowserProviderUsesCloakBrowserVersionEndpoint(t *testing.T) {
	var requestedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		assert.Equal(t, "feedcraft", r.URL.Query().Get("fingerprint"))
		_, _ = w.Write([]byte(`{"webSocketDebuggerUrl":"ws://example/devtools/browser/1"}`))
	}))
	defer server.Close()

	env := viper.New()
	env.Set("BROWSER_PROVIDER", "cloakbrowser-cdp")
	env.Set("BROWSER_ENDPOINT", server.URL+"?fingerprint=feedcraft")

	status := checkBrowserless(env, true)

	assert.Equal(t, "Browser Provider", status.Name)
	assert.Equal(t, "Healthy", status.Status)
	assert.Equal(t, "/json/version", requestedPath)
	assert.Contains(t, status.Details, "Provider: cloakbrowser-cdp")
	assert.Contains(t, status.Details, server.URL)
}
