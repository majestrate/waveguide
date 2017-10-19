package worker

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"waveguide/lib/api"
	"waveguide/lib/log"
)

func (w *Worker) UploadRequest(u *url.URL, body io.ReadCloser) *http.Request {
	return &http.Request{
		URL:    u,
		Method: "PUT",
		Body:   body,
	}
}

func (w *Worker) MkTorrentRequest(outfile string, callback *url.URL) *http.Request {
	u := w.GetNextWorkerURL()
	u.Path = fmt.Sprintf("/api/%s", api.MakeTorrent)
	q := u.Query()
	q.Add(api.ParamCallbackURL, callback.String())
	u.RawQuery = q.Encode()
	return &http.Request{
		URL: u,
		GetBody: func() (io.ReadCloser, error) {
			return os.Open(outfile)
		},
		Method: "POST",
	}
}

func (w *Worker) DoRequest(r *http.Request) error {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Errorf("error doing request: %s", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code %d", resp.StatusCode)
	}
	return err
}
