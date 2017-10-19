package util

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func RandStr(l int) string {
	buff := make([]byte, 1+(l/2))
	io.ReadFull(rand.Reader, buff)
	return hex.EncodeToString(buff)[:l]
}
