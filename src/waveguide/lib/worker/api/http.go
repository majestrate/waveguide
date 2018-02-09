package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"waveguide/lib/log"
)

func UploadRequest(u *url.URL, body io.ReadCloser) *http.Request {
	return &http.Request{
		URL:    u,
		Method: "PUT",
		Body:   body,
	}
}

func DeleteRequest(u *url.URL) *http.Request {
	return &http.Request{
		URL:    u,
		Method: "DELETE",
	}
}

func DoHTTP(r *http.Request) error {
	log.Debugf("do request %s %s", r.Method, r.URL.String())
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Errorf("error doing request: %s", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code %d", resp.StatusCode)
	}
	resp.Body.Close()
	return err
}
