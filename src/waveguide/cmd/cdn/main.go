package cdn

import (
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"waveguide/lib/log"
)

type CDNServer struct {
	rootdir string
}

func (cdn *CDNServer) HandlePUT(c *gin.Context) {
	var buff [65536]byte
	f, err := os.Create(filepath.Join(cdn.rootdir, filepath.Clean(c.Param("filename"))))
	if err == nil {
		_, err = io.CopyBuffer(f, c.Request.Body, buff[:])
		c.Request.Body.Close()
	}
	if err == nil {
		c.String(http.StatusOK, "okay")
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func Run() {
	cdn := &CDNServer{
		rootdir: "cdn_files",
	}
	os.Mkdir(cdn.rootdir, 0700)
	addr := os.Getenv("ADDR")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	router := gin.Default()
	// set up cors
	router.Use(cors.Middleware(cors.Config{
		// TODO: configure origins for CORS
		Origins:         "*",
		Methods:         "GET",
		RequestHeaders:  "Origin, Authorization, Content-Type, Range",
		ExposedHeaders:  "Range",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.PUT("/:filename", cdn.HandlePUT)
	router.Static("/", cdn.rootdir)
	sigchnl := make(chan os.Signal)
	signal.Notify(sigchnl, os.Interrupt, syscall.SIGHUP)
	go func() {
		for sig := range sigchnl {
			switch sig {
			case os.Interrupt:
				listener.Close()
			case syscall.SIGHUP:
				log.Info("SIGHUP")
				continue
			}
		}
	}()
	log.Infof("serving on %s", listener.Addr())
	http.Serve(listener, router)
}
