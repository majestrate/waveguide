package api

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

type Server struct {
	e *gin.Engine
}

func NewServer() *Server {
	return &Server{
		e: gin.Default(),
	}
}

/** setup internal state and routes */
func (s *Server) Setup() (err error) {

	return
}

func (s *Server) Serve(l net.Listener) error {
	return http.Serve(l, s.e)
}
