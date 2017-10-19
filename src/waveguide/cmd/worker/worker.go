package worker

import (
	"waveguide/lib/api"
	"waveguide/lib/database"
	"waveguide/lib/log"
	"waveguide/lib/torrent"
	"waveguide/lib/video"
)

type Worker struct {
	UploadURL string
	Encoder   video.Encoder
	Torrent   *torrent.Factory
	TempDir   string
	DB        database.Database
	API       *api.Client
}

func (w *Worker) FindWorker(method string) (api.WorkerFunc, bool) {
	log.Debugf("find worker for %s", method)
	switch method {
	case api.EncodeVideo:
		return w.ApiEncodeVideo, true
	case api.MakeTorrent:
		return w.ApiMakeTorrent, true
	default:
		return nil, false
	}
}
