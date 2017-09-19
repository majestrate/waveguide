package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"waveguide/lib/config"
	"waveguide/lib/database"
	"waveguide/lib/log"
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

	routes.DB = database.NewDatabase(conf.Frontend.DB.URL)
	err = routes.DB.Init()
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}

	// make router
	router := gin.Default()
	// set up template utils
	funcs := templates.Funcs()
	router.SetFuncMap(funcs)
	// load templates
	router.LoadHTMLGlob(filepath.Join(conf.Frontend.TemplateDir, "*.html.tmpl"))

	// static resources
	router.Static("/static", conf.Frontend.StaticDir)
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/favicon.png")
	})

	// setup routes
	router.GET("/", routes.ServeIndex)
	router.GET("/v/:VideoID/", routes.ServeVideo)
	router.GET("/u/:UserID/", routes.ServeUser)
	// run router
	router.Run()
}
