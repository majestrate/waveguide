package worker

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"waveguide/lib/log"
	"waveguide/lib/video"
	"waveguide/lib/worker/api"
)

var ErrNoFileName = errors.New("no filename provided")
var ErrNoFilePath = errors.New("no local filepath provided")
var ErrNoVideoID = errors.New("no videoid provided")

func (w *Worker) ApiEncodeVideo(r *api.Request) error {
	outfile := w.TempFileName(".mp4")
	fname := r.GetString(api.ParamFilename, "")

	if fname == "" {
		return ErrNoFileName
	}
	vidid := r.GetString(api.ParamVideoID, "")
	if vidid == "" {
		return ErrNoVideoID
	}
	infileURL, err := url.Parse(r.GetString(api.ParamFileURL, ""))
	if err != nil {
		return err
	}
	if strings.ToLower(infileURL.Scheme) != "file" {
		return ErrNoFilePath
	}
	infile := infileURL.Path
	var encode bool

	log.Infof("probing %s", infile)

	encode, err = w.Prober.VideoNeedsEncoding(infile, video.Info{
		VideoCodec: "h264",
		AudioCodec: "aac",
	})
	if err != nil {
		log.Errorf("failed to probe %s: %s", infile, err.Error())
		return err
	}

	if encode {
		log.Infof("Encoding %s as %s to %s", fname, infile, outfile)
		err = w.Encoder.EncodeFile(infile, outfile)
		os.Remove(infile)
	} else if err == nil {
		log.Infof("Video %s accepted as is", fname)
		err = os.Rename(infile, outfile)
	}
	if err == nil {
		var uploadURL *url.URL
		uploadURL, err = url.Parse(w.UploadURL)
		if err == nil {
			_, remoteFile := filepath.Split(outfile)
			uploadURL.Path = "/" + remoteFile
			var f *os.File
			f, err = os.Open(outfile)
			if err == nil {
				err = w.DoRequest(api.UploadRequest(uploadURL, f))
				f.Close()
				if err == nil {
					err = w.DB.AddWebseed(vidid, w.ToPublicCDN(uploadURL).String())
					if err == nil {
						log.Infof("make torrent for %s", outfile)
						err = w.API.Do(api.MkTorrentRequest(&url.URL{
							Path:   outfile,
							Scheme: "file",
						}, vidid, fname))
					}
				}
			}
		}
	}
	if err != nil {
		log.Errorf("failed to encode video: %s", err)
		os.Remove(outfile)
	}
	return err
}
