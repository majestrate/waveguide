package worker

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"waveguide/lib/api"
	"waveguide/lib/log"
	"waveguide/lib/model"
)

var ErrNoVideoFile = errors.New("no video file")
var ErrNoVideoInfo = errors.New("no video info")

func (w *Worker) EncodeVideo(tmpname, fname string, callback, upload_url *url.URL) {
	outfile := w.TempFileName("mp4")
	log.Infof("Encoding %s to %s", fname, outfile)
	err := w.Encoder.EncodeFile(tmpname, outfile)
	os.Remove(tmpname)
	if err == nil {
		log.Infof("Upload video file %s to %s", outfile, upload_url.String())
		err = w.DoRequest(w.MkTorrentRequest(outfile, callback))
	}
	if err != nil {
		log.Errorf("failed to encode video: %s", err)
	}
}

func (w *Worker) VideoEncode(c *gin.Context, callback *url.URL) error {
	fname := filepath.Clean(c.Query(api.ParamFilename))
	upload_url_str := c.Query(api.ParamUploadURL)
	var upload_url *url.URL
	var err error
	if upload_url_str != "" {
		upload_url, err = url.Parse(upload_url_str)
		if err != nil {
			return err
		}
	}
	f, err := w.AcquireTempFile(filepath.Ext(fname))
	if err != nil {
		return err
	}
	var buff [65536]byte
	defer f.Close()
	mediaType, params, err := mime.ParseMediaType(c.Request.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	var info *model.VideoInfo
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(c.Request.Body, params["boundary"])
		defer c.Request.Body.Close()
		for err == nil {
			var p *multipart.Part
			p, err = mr.NextPart()
			if err == io.EOF {
				err = nil
				break
			}
			switch p.FormName() {
			case api.ParamVideoFile:
				_, err = io.CopyBuffer(f, p, buff[:])
			case api.ParamVideoInfoJSON:
				info = new(model.VideoInfo)
				err = json.NewDecoder(p).Decode(info)
			}
		}
		if info == nil {
			err = ErrNoVideoInfo
		}
	} else {
		err = ErrNoVideoFile
	}
	if err != nil {
		defer os.Remove(f.Name())
		return err
	}
	go w.EncodeVideo(f.Name(), fname, callback, upload_url)
	return nil
}
