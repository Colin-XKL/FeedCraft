package util

import (
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/sirupsen/logrus"
)

// GetDefaultTargetLang returns the default target language code (e.g. "zh-CN", "en")
func GetDefaultTargetLang() string {
	envClient := GetEnvClient()
	lang := envClient.GetString("DEFAULT_TARGET_LANG")
	if lang == "" {
		return "zh-CN" // Default to Chinese
	}
	return lang
}

// GetLanguageName converts a language code to a natural language name for prompts.
// Defaults to the code itself if not found in the map.
func GetLanguageName(code string) string {
	// Simple mapping for common languages. Can be expanded.
	mapping := map[string]string{
		"zh-CN": "Simplified Chinese",
		"zh":    "Chinese",
		"en":    "English",
		"en-US": "English",
		"ja":    "Japanese",
		"ko":    "Korean",
		"fr":    "French",
		"de":    "German",
		"es":    "Spanish",
		"ru":    "Russian",
		"pt":    "Portuguese",
		"it":    "Italian",
	}

	if name, ok := mapping[code]; ok {
		return name
	}
	return code
}

// IsSameLanguage checks if the text is in the target language.
func IsSameLanguage(text string, targetLangCode string) bool {
	if strings.TrimSpace(text) == "" {
		return false
	}

	info := whatlanggo.Detect(text)
	detectedIso := info.Lang.Iso6391()

	// Normalize targetLangCode to ISO 639-1 for comparison if possible
	targetIso := targetLangCode
	if len(targetLangCode) > 2 && strings.Contains(targetLangCode, "-") {
		targetIso = strings.Split(targetLangCode, "-")[0]
	} else if len(targetLangCode) > 2 {
		// Try to match if user provided full name like "Chinese" (unlikely but possible)
		// For now assume code.
	}

	logrus.Debugf("Language detection: text_len=%d, detected=%s, target=%s", len(text), detectedIso, targetIso)

	// specific handling for Chinese because whatlanggo distinguishes Cmn (Mandarin)
	if targetIso == "zh" {
		return detectedIso == "zh" // whatlanggo returns 'zh' for Cmn? No, it returns "cmn" usually but Iso6391() returns "zh"
	}

	return detectedIso == targetIso
}
