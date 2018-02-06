package api

import (
	"encoding/json"
	"net/http"
)

const UploadEndpoint = "/upload.php"

type Backend interface {
	StorePomf(r *Request) Response
}

type Server struct {
	backend Backend
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		if r.URL.Path == UploadEndpoint {
			var resp Response
			req, err := TransformRequest(r)
			if err == nil {
				resp = s.backend.StorePomf(req)
			} else {
				resp.SetError(http.StatusInternalServerError, err)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		} else {
			http.NotFound(w, r)
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func NewServer(b Backend) *Server {
	return &Server{
		backend: b,
	}
}
