package frontend

import (
	"github.com/gin-gonic/gin"
	"waveguide/lib/config"
	"waveguide/lib/log"
)

func Run() {
	var conf config.Config

	err := conf.Load("waveguide.ini")
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	router := gin.Default()

	router.Run()
}
