package worker

import (
	"github.com/gin-gonic/gin"
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
	// TODO: do not hardcode worker url
	app.WorkerURL = conf.Frontend.WorkerURL
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

	router := gin.Default()
	router.POST("/api/:Method", app.ServeAPI)
	router.GET("/callback", app.ServeCallback)
	router.Run()
}
