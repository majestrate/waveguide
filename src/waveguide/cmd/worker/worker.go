package worker

import (
	"waveguide/lib/api"
	"waveguide/lib/config"
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
	quit      bool
}

func (w *Worker) Continue() bool {
	return !w.quit
}
func (w *Worker) Stop() {
	w.quit = true
}

func (w *Worker) Configure(conf *config.Config) error {
	return w.configure(conf, false)
}

func (w *Worker) Reconfigure(conf *config.Config) error {
	return w.configure(conf, true)
}

func (w *Worker) configure(conf *config.Config, reload bool) (err error) {
	w.Torrent, err = torrent.NewFactory()
	if err != nil {
		return
	}
	w.DB = database.NewDatabase(conf.DB.URL)
	if !reload {
		err = w.DB.Init()
	}
	if err != nil {
		return
	}
	w.UploadURL = conf.Worker.UploadURL
	w.TempDir = conf.Worker.TempDir
	w.Encoder, err = video.NewEncoder(&conf.Worker.Encoder)
	if err != nil {
		return
	}
	w.API, err = api.NewClient(&conf.MQ)
	return
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
