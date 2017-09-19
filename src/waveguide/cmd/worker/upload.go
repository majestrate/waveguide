package worker

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func (w *Worker) UploadFileRequest(u *url.URL, fname string) *http.Request {
	return &http.Request{
		URL:    u,
		Method: "PUT",
		GetBody: func() (io.ReadCloser, error) {
			return os.Open(fname)
		},
	}
}

func (w *Worker) UploadRequest(u *url.URL, body io.ReadCloser) *http.Request {
	return &http.Request{
		URL:    u,
		Method: "PUT",
		Body:   body,
	}
}

func (w *Worker) DoRequest(r *http.Request) error {
	resp, err := http.DefaultClient.Do(r)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code %d", resp.StatusCode)
	}
	return err
}
