package craft

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCommonCraftHandlerMapsUpstreamHTTPErrorToBadGateway(t *testing.T) {
	gin.SetMode(gin.TestMode)
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}))
	t.Cleanup(source.Close)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/craft/proxy?input_url="+source.URL, nil)

	CommonCraftHandlerUsingCraftOptionList(c, nil)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusBadGateway, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "http status not ok") {
		t.Fatalf("expected upstream status in body, got %s", recorder.Body.String())
	}
}

func TestCommonCraftHandlerMapsInvalidFeedToUnprocessableEntity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html><body>not a feed</body></html>"))
	}))
	t.Cleanup(source.Close)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/craft/proxy?input_url="+source.URL, nil)

	CommonCraftHandlerUsingCraftOptionList(c, nil)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusUnprocessableEntity, recorder.Body.String())
	}
	if !strings.Contains(strings.ToLower(recorder.Body.String()), "parse failed") {
		t.Fatalf("expected parse failure in body, got %s", recorder.Body.String())
	}
}

func TestHTTPStatusForCraftSourceError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  string
		want int
	}{
		{name: "upstream 403", err: "fetch failed: http status not ok: 403 Forbidden", want: http.StatusBadGateway},
		{name: "connect failed", err: "fetch failed: http get failed: dial tcp: lookup hnrss.org", want: http.StatusBadGateway},
		{name: "timeout", err: "fetch failed: http get failed: context deadline exceeded", want: http.StatusGatewayTimeout},
		{name: "parse", err: "parse failed: Failed to detect feed type", want: http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, _ := httpStatusForCraftSourceError(errString(tt.err))
			if got != tt.want {
				t.Fatalf("httpStatusForCraftSourceError(%q) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}

type errString string

func (e errString) Error() string { return string(e) }
