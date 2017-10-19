package util

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func httpOpen(u *url.URL) (io.ReadCloser, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func URLOpen(u *url.URL) (io.ReadCloser, error) {
	switch strings.ToLower(u.Scheme) {
	case "file":
		return os.Open(u.Path)
	default:
		return httpOpen(u)
	}
}
