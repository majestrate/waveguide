package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"path/filepath"
	"waveguide/lib/api"
	"waveguide/lib/config"
	"waveguide/lib/database"
	"waveguide/lib/log"
	"waveguide/lib/model"
	"waveguide/lib/templates"
)

func Run() {
	log.SetLevel("debug")
	var conf config.Config

	err := conf.Load("waveguide.ini")
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	var routes Routes

	routes.FrontendURL, err = url.Parse(conf.Frontend.FrontendURL)
	if err != nil {
		log.Fatalf("failed to parse frontend url: %s", err)
	}

	routes.DB = database.NewDatabase(conf.DB.URL)
	err = routes.DB.Init()
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}

	routes.api = api.NewClient(conf.Frontend.WorkerURL)
	routes.workerURL = conf.Frontend.WorkerURL
	// make router
	router := gin.Default()
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
	router.Run()
}
