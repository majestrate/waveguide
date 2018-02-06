package worker

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"waveguide/lib/log"
	"waveguide/lib/worker/api"
)

func (w *Worker) UploadRequest(u *url.URL, body io.ReadCloser) *http.Request {
	return &http.Request{
		URL:    u,
		Method: "PUT",
		Body:   body,
	}
}

func (w *Worker) DeleteRequest(u *url.URL) *http.Request {
	return &http.Request{
		URL:    u,
		Method: "DELETE",
	}
}

func (w *Worker) ExpireVideoRequest(vidID string) *api.Request {
	return &api.Request{
		Method: api.ExpireVideo,
		Args: map[string]interface{}{
			api.ParamVideoID: vidID,
		},
	}
}

func (w *Worker) MkTorrentRequest(infile *url.URL, vid string, filename string) *api.Request {
	return &api.Request{
		Method: api.MakeTorrent,
		Args: map[string]interface{}{
			api.ParamVideoID:  vid,
			api.ParamFilename: filename,
			api.ParamFileURL:  infile.String(),
		},
	}
}

func (w *Worker) DoRequest(r *http.Request) error {
	log.Debugf("do request %s %s", r.Method, r.URL.String())
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
