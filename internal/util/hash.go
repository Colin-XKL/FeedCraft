package util

import (
	"crypto/md5"
	"encoding/hex"
	"hash/fnv"
)

func GetTextContentHash(text string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func getMD5Hash(text string) string {
	h := md5.New()
	_, _ = h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// GetPasswordMD5Hash returns the MD5 hex digest of text, intended for password hashing flows.
func GetPasswordMD5Hash(text string) string {
	return getMD5Hash(text)
}
