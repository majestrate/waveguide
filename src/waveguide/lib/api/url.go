package api

import (
	"net/url"
	"path/filepath"
)

func (s *Server) MakeTorrentUploadURL(fname string) string {
	// TODO: multiple metainfo servers
	return s.MakeVideoUploadUrl(fname) + ".torrent"
}

func (s *Server) MakeVideoUploadUrl(fname string) string {
	fu, _ := url.Parse(s.conf.Worker.UploadURL)
	fu.Path = "/" + filepath.Base(fname)
	return fu.String()
}

func (s *Server) MakeWebseedURL(fname string) string {
	// TODO: multiple webseed servers
	fu, _ := url.Parse(s.conf.CDN.WebseedServers[0])
	fu.Path = "/" + filepath.Base(fname)
	return fu.String()
}

func (s *Server) MakePublicURL(fname string) string {
	// TODO: multiple webseed servers
	return s.MakeWebseedURL(fname) + ".torrent"
}
