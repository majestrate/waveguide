package worker

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"waveguide/lib/log"
)

func (w *Worker) EncodeVideo(tmpname, fname string, callback, upload_url *url.URL) {
	outfile := w.TempFileName("mp4")
	log.Infof("Encoding %s to %s", fname, outfile)
	err := w.Encoder.EncodeFile(tmpname, outfile)
	os.Remove(tmpname)
	if err == nil {
		log.Infof("Upload video file %s to %s", outfile, upload_url.String())
		err = w.DoRequest(w.UploadFileRequest(upload_url, outfile))
		os.Remove(outfile)
	}
	w.InformCallback(callback, err)

}

func (w *Worker) VideoEncode(c *gin.Context, callback *url.URL) error {
	fname := filepath.Clean(c.Query("filename"))
	upload_url, err := url.Parse(c.Query("upload_url"))
	if err != nil {
		return err
	}
	f, err := w.AcquireTempFile(filepath.Ext(fname))
	if err != nil {
		return err
	}
	var buff [65536]byte
	defer f.Close()
	defer c.Request.Body.Close()
	_, err = io.CopyBuffer(f, c.Request.Body, buff[:])
	if err != nil {
		defer os.Remove(f.Name())
		return err
	}
	go w.EncodeVideo(f.Name(), fname, callback, upload_url)
	return nil
}
