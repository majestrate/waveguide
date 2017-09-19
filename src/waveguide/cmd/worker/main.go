package worker

import (
	"github.com/gin-gonic/gin"
	"waveguide/lib/config"
	"waveguide/lib/log"
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

	app.TempDir = conf.Worker.TempDir
	app.Encoder, err = video.NewEncoder(&conf.Worker.Encoder)
	if err != nil {
		log.Fatalf("Error creating video encoder: %s", err)
	}

	router := gin.Default()
	router.POST("/api/:Method", app.ServeAPI)
	router.Run()
}
