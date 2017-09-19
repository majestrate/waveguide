package worker

import (
	"github.com/gin-gonic/gin"
	"net/url"
	"path/filepath"
	"waveguide/lib/util"
)

func (w *Worker) MakeTorrent(c *gin.Context, u *url.URL) error {
	fname := filepath.Clean(c.Query("filename"))
	uploadURL, err := url.Parse(c.Query("upload_url"))
	if err == nil {
		torrent := new(util.Buffer)
		go func(c *gin.Context) {
			err = w.Torrent.MakeSingle(fname, c.Request.Body, torrent)
			c.Request.Body.Close()
			if err == nil {
				err = w.DoRequest(w.UploadRequest(uploadURL, torrent))
				w.InformCallback(u, err)

			}
		}(c.Copy())
	}
	return err
}
