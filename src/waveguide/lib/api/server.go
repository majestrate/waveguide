package api

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"waveguide/lib/config"
	"waveguide/lib/oauth"
	"waveguide/lib/streaming"
	"waveguide/lib/torrent"
)

type Server struct {
	e       *gin.Engine
	conf    config.Config
	oauth   *oauth.Client
	torrent *torrent.Factory
	ctx     *streaming.Context
}

func (s *Server) Configure(conf config.Config) (err error) {
	if s.e == nil {
		s.setupRoutes()
		s.ctx = streaming.NewContext()
		err = s.reconfigure(conf, true)
	} else {
		err = s.reconfigure(conf, false)
	}
	return
}

func (s *Server) reconfigure(conf config.Config, fresh bool) (err error) {
	s.conf = conf
	s.oauth = oauth.NewClient(s.conf.OAuth)
	s.torrent, err = torrent.NewFactory(&s.conf.Worker.Torrent)
	return
}

func (s *Server) setupRoutes() {
	s.e = gin.Default()
	s.e.POST("/api/v1/stream/publish", s.APIStreamPublish)
	s.e.POST("/api/v1/stream/join", s.APIStreamJoin)
	s.e.POST("/api/v1/stream/part", s.APIStreamPart)
	s.e.POST("/api/v1/stream/done", s.APIStreamDone)
	s.e.POST("/api/v1/stream/segment", s.APIStreamSegment)
	s.e.GET("/api/v1/stream/", s.APIStreamInfo)
}

func (s *Server) Serve(l net.Listener) error {
	return http.Serve(l, s.e)
}
