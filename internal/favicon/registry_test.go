package favicon

import (
	"FeedCraft/internal/config"
	"fmt"
	"sync"
	"testing"
)

func TestBuildURLUsesDefaultProvider(t *testing.T) {
	resetRegistryForTest(t)

	got, providerID := BuildURL("", "https://example.com/blog", 64)

	const want = "https://t0.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fexample.com&size=64"
	if got != want {
		t.Fatalf("BuildURL() = %q, want %q", got, want)
	}
	if providerID != ProviderGstaticCN {
		t.Fatalf("provider ID = %q, want %q", providerID, ProviderGstaticCN)
	}
}

func TestReplaceSupportsCustomTemplate(t *testing.T) {
	resetRegistryForTest(t)
	settings := config.FaviconSettings{
		DefaultProviderID: "internal",
		CustomProviders: []config.FaviconProviderConfig{
			{
				ID:          "internal",
				Name:        "Internal favicon",
				URLTemplate: "https://icons.example.test/favicon?origin={origin_query}&host={host}&size={size}",
				Enabled:     true,
			},
		},
	}

	if err := Replace(settings); err != nil {
		t.Fatalf("Replace() error = %v", err)
	}

	got, providerID := BuildURL("", "https://www.example.com/path?q=1", 96)
	const want = "https://icons.example.test/favicon?origin=https%3A%2F%2Fwww.example.com&host=www.example.com&size=96"
	if got != want {
		t.Fatalf("BuildURL() = %q, want %q", got, want)
	}
	if providerID != "internal" {
		t.Fatalf("provider ID = %q, want internal", providerID)
	}
}

func TestRecipeProviderOverridesGlobalDefault(t *testing.T) {
	resetRegistryForTest(t)
	if err := Replace(config.FaviconSettings{DefaultProviderID: ProviderGoogle}); err != nil {
		t.Fatalf("Replace() error = %v", err)
	}

	got, providerID := BuildURL(ProviderDuckDuckGo, "https://example.com/path", 64)

	if got != "https://icons.duckduckgo.com/ip3/example.com.ico" {
		t.Fatalf("BuildURL() = %q", got)
	}
	if providerID != ProviderDuckDuckGo {
		t.Fatalf("provider ID = %q, want %q", providerID, ProviderDuckDuckGo)
	}
}

func TestUnknownProviderFallsBackToConfiguredDefault(t *testing.T) {
	resetRegistryForTest(t)
	if err := Replace(config.FaviconSettings{DefaultProviderID: ProviderYandex}); err != nil {
		t.Fatalf("Replace() error = %v", err)
	}

	got, providerID := BuildURL("missing", "https://example.com/path", 64)

	if got != "https://favicon.yandex.net/favicon/example.com" {
		t.Fatalf("BuildURL() = %q", got)
	}
	if providerID != ProviderYandex {
		t.Fatalf("provider ID = %q, want %q", providerID, ProviderYandex)
	}
}

func TestReplaceRejectsInvalidCustomProvidersAndKeepsSnapshot(t *testing.T) {
	resetRegistryForTest(t)
	before, beforeProvider := BuildURL("", "https://example.com", 64)

	tests := []config.FaviconProviderConfig{
		{ID: ProviderGoogle, Name: "Reserved", URLTemplate: "https://icons.example.test/{host}", Enabled: true},
		{ID: "bad id", Name: "Bad ID", URLTemplate: "https://icons.example.test/{host}", Enabled: true},
		{ID: "unknown_placeholder", Name: "Bad placeholder", URLTemplate: "https://icons.example.test/{domain}", Enabled: true},
		{ID: "http_provider", Name: "HTTP", URLTemplate: "http://icons.example.test/{host}", Enabled: true},
		{ID: "no_target", Name: "No target", URLTemplate: "https://icons.example.test/icon?size={size}", Enabled: true},
	}

	for _, provider := range tests {
		t.Run(provider.ID, func(t *testing.T) {
			err := Replace(config.FaviconSettings{
				DefaultProviderID: provider.ID,
				CustomProviders:   []config.FaviconProviderConfig{provider},
			})
			if err == nil {
				t.Fatal("Replace() error = nil, want validation error")
			}
			after, afterProvider := BuildURL("", "https://example.com", 64)
			if after != before || afterProvider != beforeProvider {
				t.Fatalf("snapshot changed after failed Replace(): (%q, %q)", after, afterProvider)
			}
		})
	}
}

func TestDisabledProviderFallsBackToDefault(t *testing.T) {
	resetRegistryForTest(t)
	if err := Replace(config.FaviconSettings{
		DefaultProviderID: ProviderGstaticCN,
		CustomProviders: []config.FaviconProviderConfig{
			{
				ID:          "disabled",
				Name:        "Disabled",
				URLTemplate: "https://icons.example.test/{host}",
				Enabled:     false,
			},
		},
	}); err != nil {
		t.Fatalf("Replace() error = %v", err)
	}

	_, providerID := BuildURL("disabled", "https://example.com", 64)
	if providerID != ProviderGstaticCN {
		t.Fatalf("provider ID = %q, want %q", providerID, ProviderGstaticCN)
	}
}

func TestCustomProviderRuntimeHTTPOutputFallsBackToGstaticCN(t *testing.T) {
	resetRegistryForTest(t)
	if err := Replace(config.FaviconSettings{
		DefaultProviderID: "raw_origin",
		CustomProviders: []config.FaviconProviderConfig{
			{
				ID:          "raw_origin",
				Name:        "Raw origin",
				URLTemplate: "{origin}/favicon.ico",
				Enabled:     true,
			},
		},
	}); err != nil {
		t.Fatalf("Replace() error = %v", err)
	}

	got, providerID := BuildURL("", "http://example.com/path", 64)

	const want = "https://t0.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=http%3A%2F%2Fexample.com&size=64"
	if got != want {
		t.Fatalf("BuildURL() = %q, want %q", got, want)
	}
	if providerID != ProviderGstaticCN {
		t.Fatalf("provider ID = %q, want %q", providerID, ProviderGstaticCN)
	}
}

func TestRegistryConcurrentBuildAndReplace(t *testing.T) {
	resetRegistryForTest(t)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				got, _ := BuildURL("", fmt.Sprintf("https://example%d.com/path", j), 64)
				if got == "" {
					t.Errorf("BuildURL() returned empty URL")
					return
				}
			}
		}()
	}
	for i := 0; i < 50; i++ {
		defaultProvider := ProviderGstaticCN
		if i%2 == 1 {
			defaultProvider = ProviderYandex
		}
		if err := Replace(config.FaviconSettings{DefaultProviderID: defaultProvider}); err != nil {
			t.Fatalf("Replace() error = %v", err)
		}
	}
	wg.Wait()
}

func resetRegistryForTest(t *testing.T) {
	t.Helper()
	previous := Settings()
	t.Cleanup(func() {
		if err := Replace(previous); err != nil {
			t.Fatalf("restore registry: %v", err)
		}
	})
	if err := Replace(DefaultSettings()); err != nil {
		t.Fatalf("reset registry: %v", err)
	}
}
