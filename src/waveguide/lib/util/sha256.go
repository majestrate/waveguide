package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(str string) string {
	d := sha256.Sum256([]byte(str))
	return hex.EncodeToString(d[:])
}
