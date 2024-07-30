package util

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
