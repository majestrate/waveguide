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
	router.HEAD("/", func(c *gin.Context) {
		c.String(200, "")
	})
	router.GET(fmt.Sprintf("%s/:id/", model.VideoURLBase), routes.ServeVideo)
	router.HEAD(fmt.Sprintf("%s/:id/", model.VideoURLBase), func(c *gin.Context) {
		c.String(200, "")
	})

	router.GET("/u/:Username/", routes.ServeUser)
	router.GET("/u/:Username/videos.atom", routes.ServeUserVideosFeed)

	router.GET("/oauth/redirect_uri", routes.HandleOAuthRedirect)

	apiV1 := router.Group("/wg-api/v1")
	{
		apiV1.GET("/online", routes.ApiStreamsOnline)
		/*
			apiV1.POST("/login", routes.ApiLogin)
			apiV1.POST("/register", routes.ApiRegister)
		*/
		apiV1.GET("/stream", routes.ApiStreamMagnets)
		apiV1.GET("/comments", routes.ApiStreamComments)
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
	router.GET("/logout/", routes.ApiLogout)

	router.GET("/stream/", routes.ServeStream)
	router.GET("/watch/", routes.ServeWatch)
	// chat callback url
	router.StaticFile("/chat/", filepath.Join(conf.Frontend.StaticDir, "chat.html"))
	router.GET("/register/", routes.ServeRegister)

	// pomf api
	if routes.Pomf != nil {
		router.Any("/upload.php", gin.WrapH(routes.Pomf))
	}

}
