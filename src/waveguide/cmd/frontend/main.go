package frontend

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"waveguide/lib/config"
	"waveguide/lib/log"
	"waveguide/lib/model"
	"waveguide/lib/templates"
	"waveguide/lib/util"
)

func Run() {
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
	router.Use(util.CORSMiddleware())

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
	router.GET("/u/:Username/", routes.ServeUser)
	router.GET("/u/:Username/videos.atom", routes.ServeUserVideosFeed)

	apiV1 := router.Group("/wg-api/v1")
	{
		apiV1.POST("/login", routes.ApiLogin)
		apiV1.POST("/register", routes.ApiRegister)
		authed := apiV1.Group("/authed")
		authed.Use(routes.ApiAuthMiddleware())
		{
			authed.POST("/upload", routes.ApiUpload)
			authed.POST("/comment", routes.ApiComment)
		}
	}
	router.GET("/captcha/:f", gin.WrapH(captcha.Server(500, 200)))

	router.GET("/upload/", routes.ServeUpload)
	router.GET("/login/", routes.ServeLogin)

	router.GET("/register/", routes.ServeRegister).Use(RequiresCaptchaMiddleware())

	// chat callback url
	router.StaticFile("/chat/", filepath.Join(conf.Frontend.StaticDir, "chat.html"))
	// run router
	log.Infof("running on %s", listener.Addr())
	http.Serve(listener, router)
}
