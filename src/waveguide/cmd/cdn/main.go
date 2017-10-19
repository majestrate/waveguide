package cdn

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	router := gin.Default()
	router.PUT("/:filename", cdn.HandlePUT)
	router.Run()
}
