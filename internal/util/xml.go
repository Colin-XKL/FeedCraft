package util

import (
	"bytes"
	"strings"
)

// StripInvalidXMLChars removes characters that are illegal in XML 1.0 documents.
// XML 1.0 Char ::= #x9 | #xA | #xD | [#x20-#xD7FF] | [#xE000-#xFFFD] | [#x10000-#x10FFFF]
func StripInvalidXMLChars(s string) string {
	if !strings.ContainsFunc(s, isInvalidXMLRune) {
		return s
	}
	return strings.Map(func(r rune) rune {
		if isXMLChar(r) {
			return r
		}
		return -1
	}, s)
}

// StripInvalidXMLBytes removes XML 1.0 illegal characters from a byte slice.
func StripInvalidXMLBytes(data []byte) []byte {
	if len(data) == 0 || !bytes.ContainsFunc(data, isInvalidXMLRune) {
		return data
	}
	return []byte(StripInvalidXMLChars(string(data)))
}

func isInvalidXMLRune(r rune) bool {
	return !isXMLChar(r)
}

func isXMLChar(r rune) bool {
	switch {
	case r == 0x09 || r == 0x0A || r == 0x0D:
		return true
	case r >= 0x20 && r <= 0xD7FF:
		return true
	case r >= 0xE000 && r <= 0xFFFD:
		return true
	case r >= 0x10000 && r <= 0x10FFFF:
		return true
	default:
		return false
	}
}
