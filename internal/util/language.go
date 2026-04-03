package util

import (
	"regexp"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var (
	htmlTagRegex = regexp.MustCompile(`(?i)<[^>]*>`)
	urlRegex     = regexp.MustCompile(`(?i)https?://[^\s]+`)
	uuidRegex    = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	enWordRegex  = regexp.MustCompile(`(?i)\b[a-z]+\b`)
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

	// Clean text from noise (HTML tags, URLs, UUIDs) that distorts trigram-based language detection
	cleanText := uuidRegex.ReplaceAllString(text, " ")
	cleanText = htmlTagRegex.ReplaceAllString(cleanText, " ")
	cleanText = urlRegex.ReplaceAllString(cleanText, " ")

	info := whatlanggo.Detect(cleanText)
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

		return checkLanguageMatchFallback(cleanText, detectedIso, targetIso)
	}

	base, _ := targetTag.Base()
	targetIso := base.String()

	logrus.Debugf("Language detection: text_len=%d, detected=%s, target=%s", len(text), detectedIso, targetIso)

	return checkLanguageMatchFallback(cleanText, detectedIso, targetIso)
}

// checkLanguageMatchFallback applies additional CJK-specific language detection fallback.
// English words in short texts can heavily sway the whatlanggo trigram models, causing
// Japanese/Chinese texts to be detected as English, Danish, Portuguese, etc.
func checkLanguageMatchFallback(cleanText string, detectedIso string, targetIso string) bool {
	if detectedIso == targetIso {
		return true
	}

	// For CJK languages, strip English alphabetical words and try detection again
	if targetIso == "ja" || targetIso == "zh" || targetIso == "ko" {
		cleanTextCJK := enWordRegex.ReplaceAllString(cleanText, " ")
		if strings.TrimSpace(cleanTextCJK) != "" {
			infoCJK := whatlanggo.Detect(cleanTextCJK)
			detectedIsoCJK := infoCJK.Lang.Iso6391()

			// If whatlanggo successfully detects the target CJK language after stripping English text
			if detectedIsoCJK == targetIso {
				logrus.Debugf("Language detection fallback triggered: original=%s, corrected=%s", detectedIso, detectedIsoCJK)
				return true
			}
		}
	}

	// whatlanggo returns "zh" for Mandarin (Cmn) when calling Iso6391()
	return detectedIso == targetIso
}
