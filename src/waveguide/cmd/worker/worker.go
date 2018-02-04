package worker

import (
	"net/url"
	"waveguide/lib/config"
	"waveguide/lib/database"
	"waveguide/lib/log"
	"waveguide/lib/torrent"
	"waveguide/lib/video"
	"waveguide/lib/worker/api"
)

type Worker struct {
	UploadURL string
	Encoder   video.Encoder
	Prober    video.Prober
	Torrent   *torrent.Factory
	TempDir   string
	DB        database.Database
	API       *api.Client
	quit      bool
	CDN       config.CDNConfig
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
	w.CDN = conf.CDN
	w.Torrent, err = torrent.NewFactory(&conf.Worker.Torrent)
	if err != nil {
		return
	}
	if w.DB != nil {
		w.DB.Close()
	}
	w.DB = database.NewDatabase(conf.DB.URL)
	err = w.DB.Init()
	if err != nil {
		return
	}
	w.UploadURL = conf.Worker.UploadURL
	w.TempDir = conf.Worker.TempDir
	w.Encoder, err = video.NewEncoder(&conf.Worker.Encoder)
	if err != nil {
		log.Fatalf("failed to create encoder: %s", err)
		return
	}
	w.Prober, err = video.NewProber(&conf.Worker.Encoder)
	if err != nil {
		log.Fatalf("failed to create prober: %s", err)
		return
	}

	if w.API != nil {
		w.API.Close()
	}

	w.API, err = api.NewClient(&conf.MQ)
	return
}

func (w *Worker) ToPublicCDN(u *url.URL) *url.URL {
	// TODO: implement this better
	cdn, _ := url.Parse(w.CDN.WebseedServers[0])
	cdn.Path += u.Path[1:]
	return cdn
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
