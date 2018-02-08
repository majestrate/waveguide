package api

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"waveguide/lib/config"
	"waveguide/lib/log"
	"waveguide/lib/oauth"
	"waveguide/lib/streaming"
	"waveguide/lib/torrent"
	"waveguide/lib/video"
)

type Server struct {
	e       *gin.Engine
	conf    config.Config
	oauth   *oauth.Client
	torrent *torrent.Factory
	ctx     *streaming.Context
	encoder video.Encoder
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
	s.encoder, err = video.NewEncoder(&s.conf.Worker.Encoder)
	if err != nil {
		log.Fatalf("failed to create encoder: %s", err)
		return
	}
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
	s.e.GET("/api/v1/stream/info/:key", s.APIStreamInfo)
	s.e.GET("/api/v1/streams/", s.APIListStreams)
}

func (s *Server) Serve(l net.Listener) error {
	return http.Serve(l, s.e)
}
