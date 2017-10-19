package util

import (
	"net/url"
	"os"
	"strings"
)

func RemoveURL(u *url.URL) error {
	switch strings.ToLower(u.Scheme) {
	case "file":
		return os.Remove(u.Path)
	default:
		return nil
	}
}
