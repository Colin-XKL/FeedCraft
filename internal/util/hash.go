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

func GetPasswordMD5Hash(text string) string {
	h := md5.New()
	_, _ = h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
