package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"waveguide/lib/config"
	"waveguide/lib/log"
	"waveguide/lib/model"
	"waveguide/lib/templates"
)

func Run() {
	log.SetLevel("debug")
	var conf config.Config

	const configFname = "waveguide.ini"

	err := conf.Load(configFname)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	routes := new(Routes)

	err = routes.Configure(&conf)
	if err != nil {
		log.Fatalf("failed to configure frontend: %s", err)
	}

	// make net listener
	var listener net.Listener
	listener, err = net.Listen("tcp", conf.Frontend.Addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	// make router
	router := gin.Default()

	// set up cors
	router.Use(cors.Middleware(cors.Config{
		// TODO: configure origins for CORS
		Origins:         "*",
		Methods:         "GET",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	sigchnl := make(chan os.Signal)
	signal.Notify(sigchnl, os.Interrupt, syscall.SIGHUP)
	go func(chnl chan os.Signal) {
		for sig := range chnl {
			switch sig {
			case os.Interrupt:
				listener.Close()
				routes.Close()
			case syscall.SIGHUP:
				log.Info("SIGHUP")
				err := conf.Load(configFname)
				if err == nil {
					log.Info("reconfiguring")
					err = routes.Reconfigure(&conf)
				}
				if err != nil {
					log.Errorf("failed to reconfigure: %s", err)
				}
			}
		}
	}(sigchnl)

	// set up template utils
	funcs := templates.Funcs()
	router.SetFuncMap(funcs)
	// load templates
	router.LoadHTMLGlob(filepath.Join(conf.Frontend.TemplateDir, "*.html"))

	// static resources
	router.Static("/static", conf.Frontend.StaticDir)
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/favicon.png")
	})

	// setup routes
	router.GET("/", routes.ServeIndex)
	router.GET(fmt.Sprintf("%s/:VideoID/", model.VideoURLBase), routes.ServeVideo)
	router.GET("/u/:UserID/", routes.ServeUser)
	router.GET("/upload/", routes.ServeUpload)
	router.POST("/upload/", routes.HandleUpload)
	// run router
	log.Infof("running on %s", listener.Addr())
	http.Serve(listener, router)
}
