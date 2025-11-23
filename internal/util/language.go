package util

import (
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// GetDefaultTargetLang returns the default target language code (e.g. "zh-CN", "en")
// It validates the input and falls back to "zh-CN" if invalid.
func GetDefaultTargetLang() string {
	envClient := GetEnvClient()
	lang := envClient.GetString("DEFAULT_TARGET_LANG")
	if lang == "" {
		return "zh-CN" // Default to Chinese
	}

	// Validate the language code
	_, err := language.Parse(lang)
	if err != nil {
		logrus.Warnf("Invalid DEFAULT_TARGET_LANG '%s': %v. Falling back to zh-CN.", lang, err)
		return "zh-CN"
	}

	return lang
}

// GetLanguageName converts a language code to a natural language name for prompts (in English).
// It uses golang.org/x/text/language/display to get the name.
func GetLanguageName(code string) string {
	tag, err := language.Parse(code)
	if err != nil {
		return code
	}

	// Get the English name of the language
	name := display.English.Languages().Name(tag)
	if name == "" {
		return code
	}

	// For Chinese, specific handling to return "Simplified Chinese" if applicable,
	// or just "Chinese" if generic.
	// display package handles this:
	// zh-CN -> Chinese (Simplified)
	// zh -> Chinese

	return name
}

// IsSameLanguage checks if the text is in the target language.
func IsSameLanguage(text string, targetLangCode string) bool {
	if strings.TrimSpace(text) == "" {
		return false
	}

	info := whatlanggo.Detect(text)
	detectedIso := info.Lang.Iso6391()

	// Normalize targetLangCode to ISO 639-1 for comparison
	// We use language.Parse to handle normalization robustly
	targetTag, err := language.Parse(targetLangCode)
	if err != nil {
		// Fallback to simple split if parse fails (though GetDefaultTargetLang should have caught this)
		targetIso := targetLangCode
		if len(targetLangCode) > 2 && strings.Contains(targetLangCode, "-") {
			targetIso = strings.Split(targetLangCode, "-")[0]
		}
		return detectedIso == targetIso
	}

	base, _ := targetTag.Base()
	targetIso := base.String()

	logrus.Debugf("Language detection: text_len=%d, detected=%s, target=%s", len(text), detectedIso, targetIso)

	// whatlanggo returns "zh" for Mandarin (Cmn) when calling Iso6391()
	return detectedIso == targetIso
}
