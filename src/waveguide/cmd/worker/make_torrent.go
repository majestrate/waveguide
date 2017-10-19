package worker

import (
	"net/url"
	"path/filepath"
	"waveguide/lib/api"
	"waveguide/lib/log"
	"waveguide/lib/util"
)

func (w *Worker) ApiMakeTorrent(r *api.Request) error {
	fileURL, err := url.Parse(r.GetString(api.ParamFileURL, ""))
	if err != nil {
		return err
	}
	filename := r.GetString(api.ParamFilename, "")
	uploadURL, err := url.Parse(w.UploadURL)
	if err != nil {
		return err
	}
	_, remoteFile := filepath.Split(fileURL.Path)
	uploadURL.Path = "/" + remoteFile + ".torrent"
	f, err := util.URLOpen(fileURL)
	if err != nil {
		return err
	}
	log.Debugf("make torrent %s for file at %s upload to %s", filename, fileURL.String(), uploadURL.String())
	torrent := new(util.Buffer)
	err = w.Torrent.MakeSingle(filename, f, torrent)
	f.Close()
	if err == nil {
		err = w.DoRequest(w.UploadRequest(uploadURL, torrent))
	}
	return err
}
