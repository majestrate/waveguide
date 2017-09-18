package frontend

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"waveguide/lib/config"
	"waveguide/lib/database"
	"waveguide/lib/log"
	"waveguide/lib/templates"
)

func Run() {
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
	// setup routes
	router.GET("/", routes.Index)
	// run router
	router.Run()
}
