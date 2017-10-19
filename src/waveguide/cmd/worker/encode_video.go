package worker

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"waveguide/lib/api"
	"waveguide/lib/log"
)

var ErrNoFileName = errors.New("no filename provided")
var ErrNoFilePath = errors.New("no local filepath provided")

func (w *Worker) ApiEncodeVideo(r *api.Request) error {
	outfile := w.TempFileName(".mp4")
	fname := r.GetString(api.ParamFilename, "")

	if fname == "" {
		return ErrNoFileName
	}

	infileURL, err := url.Parse(r.GetString(api.ParamFileURL, ""))
	if err != nil {
		return err
	}
	if strings.ToLower(infileURL.Scheme) != "file" {
		return ErrNoFilePath
	}
	infile := infileURL.Path
	log.Infof("Encoding %s as %s to %s", fname, infile, outfile)
	err = w.Encoder.EncodeFile(infile, outfile)
	os.Remove(infile)
	if err == nil {
		var uploadURL *url.URL
		uploadURL, err = url.Parse(w.UploadURL)
		if err == nil {
			_, remoteFile := filepath.Split(outfile)
			uploadURL.Path = "/" + remoteFile
			var f *os.File
			f, err = os.Open(outfile)
			if err == nil {
				err = w.DoRequest(w.UploadRequest(uploadURL, f))
				f.Close()
				if err == nil {
					log.Infof("make torrent for %s", outfile)
					err = w.API.Do(w.MkTorrentRequest(&url.URL{
						Path:   outfile,
						Scheme: "file",
					}, fname))
				}
			}
		}
	}
	if err != nil {
		log.Errorf("failed to encode video: %s", err)
	}
	return err
}
