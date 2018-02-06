package worker

import (
	"net/url"
	"path/filepath"
	"waveguide/lib/expire"
	"waveguide/lib/log"
	"waveguide/lib/util"
	"waveguide/lib/worker/api"
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
		if err == nil {
			vidid := r.GetString(api.ParamVideoID, "")
			if vidid != "" {
				err = w.DB.SetVideoMetaInfo(vidid, w.ToPublicCDN(uploadURL).String())
			}
		}
	}
	util.RemoveURL(fileURL)
	shouldExpire := util.RandBool(0.25)
	if shouldExpire || true {
		vids, err := w.DB.GetExpiredVideos(expire.DefaultCapacity)
		if err == nil {
			for _, vid := range vids {
				w.API.Do(w.ExpireVideoRequest(vid.VideoID))
			}
		} else {
			log.Errorf("failed to get expired videos: %s", err.Error())
		}
	}
	return err
}
