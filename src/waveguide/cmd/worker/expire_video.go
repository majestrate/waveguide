package worker

import (
	"net/url"
	"waveguide/lib/model"
	"waveguide/lib/worker/api"
)

func (w *Worker) ApiExpireVideo(r *api.Request) (err error) {
	id := r.GetString(api.ParamVideoID, "")
	if id == "" {
		err = ErrNoVideoID
	} else {
		var u, cdn *url.URL
		cdn, err = url.Parse(w.UploadURL)
		var info *model.VideoInfo
		info, err = w.DB.GetVideoInfo(id)
		if err == nil {
			err = w.DB.DeleteVideo(id)
		}
		if info != nil {
			for _, webseed := range info.WebSeeds {
				u, err = url.Parse(webseed)
				if u != nil {
					if u.Host != cdn.Host {
						u.Host = cdn.Host
						u.Scheme = "http"
						w.DoRequest(api.DeleteRequest(u))
					}
				}
			}

			u, err = url.Parse(info.TorrentURL)
			if u != nil {
				u.Host = cdn.Host
				u.Scheme = "http"
				w.DoRequest(api.DeleteRequest(u))
			}
		}
	}
	return
}
