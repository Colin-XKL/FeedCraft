package util

import (
	"testing"
)

func TestGetLanguageName(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"zh-CN", "Simplified Chinese"},
		{"zh", "Chinese"},
		{"en", "English"},
		{"en-US", "English"},
		{"unknown", "unknown"},
	}

	for _, test := range tests {
		result := GetLanguageName(test.code)
		if result != test.expected {
			t.Errorf("GetLanguageName(%q) = %q; want %q", test.code, result, test.expected)
		}
	}
}

func TestIsSameLanguage(t *testing.T) {
	tests := []struct {
		text           string
		targetLangCode string
		expected       bool
	}{
		{"This is an English text.", "en", true},
		{"This is an English text.", "zh", false},
		{"这是一个中文句子。", "zh-CN", true},
		{"这是一个中文句子。", "zh", true},
		{"这是一个中文句子。", "en", false},
		// Mixed content, predominantly Chinese
		{"Hello 这是一个中文句子。", "zh", true},
		// Short text might be hard to detect, but let's see
		{"你好", "zh", true},
	}

	for _, test := range tests {
		result := IsSameLanguage(test.text, test.targetLangCode)
		if result != test.expected {
			t.Errorf("IsSameLanguage(%q, %q) = %v; want %v", test.text, test.targetLangCode, result, test.expected)
		}
	}
}
