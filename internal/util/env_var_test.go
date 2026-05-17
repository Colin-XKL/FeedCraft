package util

import (
	"testing"
)

func TestDefaultFeedUserAgent(t *testing.T) {
	t.Setenv("FC_HTTP_USER_AGENT_FEED", "")
	if got := DefaultFeedUserAgent(); got != defaultFeedUserAgent {
		t.Fatalf("expected default feed user agent %q, got %q", defaultFeedUserAgent, got)
	}

	t.Setenv("FC_HTTP_USER_AGENT_FEED", "CustomFeedUA/1.0")
	if got := DefaultFeedUserAgent(); got != "CustomFeedUA/1.0" {
		t.Fatalf("expected custom feed user agent, got %q", got)
	}
}

func TestDefaultHTMLUserAgent(t *testing.T) {
	t.Setenv("FC_HTTP_USER_AGENT_HTML", "")
	if got := DefaultHTMLUserAgent(); got != htmlDefaultUserAgent {
		t.Fatalf("expected default html user agent %q, got %q", htmlDefaultUserAgent, got)
	}

	t.Setenv("FC_HTTP_USER_AGENT_HTML", "CustomHTMLUA/2.0")
	if got := DefaultHTMLUserAgent(); got != "CustomHTMLUA/2.0" {
		t.Fatalf("expected custom html user agent, got %q", got)
	}
}
