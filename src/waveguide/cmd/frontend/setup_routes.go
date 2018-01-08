package frontend

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"waveguide/lib/config"
	"waveguide/lib/model"
)

func (routes *Routes) SetupRoutes(router *gin.Engine, conf *config.Config) {
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
		apiV1.GET("/stream/:UserID", routes.ApiStreamMagnets)
		authed := apiV1.Group("/authed")
		authed.Use(routes.ApiAuthMiddleware())
		{
			authed.POST("/upload", routes.ApiUpload)
			authed.POST("/comment", routes.ApiComment)
			authed.POST("/stream-update", routes.ApiStreamUpdate)
		}

	}
	router.GET("/captcha/:f", gin.WrapH(captcha.Server(500, 200)))

	router.GET("/upload/", routes.ServeUpload)
	router.GET("/login/", routes.ServeLogin)
	router.GET("/stream/", routes.ServeStream)
	// chat callback url
	router.StaticFile("/chat/", filepath.Join(conf.Frontend.StaticDir, "chat.html"))

	router.GET("/register/", routes.ServeRegister).Use(RequiresCaptchaMiddleware())

}