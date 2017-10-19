package worker

import (
	"waveguide/lib/api"
	"waveguide/lib/config"
	"waveguide/lib/database"
	"waveguide/lib/log"
	"waveguide/lib/torrent"
	"waveguide/lib/video"
)

func Run() {
	log.SetLevel("debug")
	var app Worker
	var conf config.Config
	err := conf.Load("waveguide.ini")
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	app.UploadURL = conf.Worker.UploadURL
	app.TempDir = conf.Worker.TempDir
	app.Encoder, err = video.NewEncoder(&conf.Worker.Encoder)
	if err != nil {
		log.Fatalf("Error creating video encoder: %s", err)
	}

	app.Torrent, err = torrent.NewFactory()
	if err != nil {
		log.Fatalf("failed to create torrent factory: %s", err)
	}

	app.DB = database.NewDatabase(conf.DB.URL)
	err = app.DB.Init()
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err)
	}

	app.API = api.NewClient(&conf.MQ)

	w := api.NewWorker(&conf.MQ, app.FindWorker)
	err = w.Run()
	if err != nil {
		log.Fatalf("worker died: %s", err)
	}
}
