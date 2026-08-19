package util

import "testing"

func TestBuildFaviconURLDefaultUsesGstaticCN(t *testing.T) {
	got := BuildFaviconURL("https://example.com/blog")
	const want = "https://t0.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fexample.com&size=64"
	if got != want {
		t.Fatalf("BuildFaviconURL() = %q, want %q", got, want)
	}
}

func TestBuildFaviconURLWithProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider FaviconProvider
		pageURL  string
		want     string
	}{
		{
			name:     "gstatic cn",
			provider: FaviconProviderGstaticCN,
			pageURL:  "https://example.com/path",
			want:     "https://t0.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fexample.com&size=64",
		},
		{
			name:     "google",
			provider: FaviconProviderGoogle,
			pageURL:  "https://example.com/path",
			want:     "https://www.google.com/s2/favicons?domain_url=https%3A%2F%2Fexample.com&sz=64",
		},
		{
			name:     "duckduckgo",
			provider: FaviconProviderDuckDuckGo,
			pageURL:  "https://www.example.com/path",
			want:     "https://icons.duckduckgo.com/ip3/www.example.com.ico",
		},
		{
			name:     "yandex",
			provider: FaviconProviderYandex,
			pageURL:  "https://example.com/path",
			want:     "https://favicon.yandex.net/favicon/example.com",
		},
		{
			name:     "unknown falls back to default",
			provider: FaviconProvider("unknown"),
			pageURL:  "https://example.com/path",
			want:     "https://t0.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fexample.com&size=64",
		},
		{
			name:     "empty page url",
			provider: FaviconProviderGstaticCN,
			pageURL:  "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildFaviconURLWithProvider(tt.provider, tt.pageURL)
			if got != tt.want {
				t.Fatalf("BuildFaviconURLWithProvider(%q, %q) = %q, want %q", tt.provider, tt.pageURL, got, tt.want)
			}
		})
	}
}
