package api

import (
	"io"
	"net/http"
)

func (c *Client) UploadVideo(callbackURL, fname string, body io.ReadCloser) error {
	u := c.createURL(VideoEncode, callbackURL)
	r := &http.Request{
		URL:  u,
		Body: body,
	}
	return c.doRequest(r)
}
