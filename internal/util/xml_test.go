package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripInvalidXMLChars(t *testing.T) {
	assert.Equal(t, "helloworld", StripInvalidXMLChars("hello\x05world"))
	assert.Equal(t, "keep\tnew\nline\r", StripInvalidXMLChars("keep\tnew\nline\r"))
	assert.Equal(t, "ok", StripInvalidXMLChars("ok"))
	assert.Equal(t, "", StripInvalidXMLChars("\x00\x01\x05"))
}

func TestStripInvalidXMLBytes(t *testing.T) {
	input := []byte("hello\x05world")
	assert.Equal(t, []byte("helloworld"), StripInvalidXMLBytes(input))
	unchanged := []byte("already valid")
	assert.Equal(t, unchanged, StripInvalidXMLBytes(unchanged))
}
