package frontend

import (
	"io"
	"net/http"
	"waveguide/lib/model"
	"waveguide/lib/pomf/api"
)

func (r *Routes) StorePomf(req *api.Request) (resp api.Response) {
	var err error
	err = req.WalkFiles(func(filename string, body io.Reader) (e error) {
		var url string
		// ANON User
		info := NewUpload(new(model.UserInfo))
		info.Title = filename
		url, e = r.UploadVideoFile(filename, body, info)
		if e == nil {
			resp.AddFile(api.Resource{
				Name: filename,
				URL:  url,
			})
		}
		return
	})
	if err != nil {
		resp.SetError(http.StatusInternalServerError, err)
	}
	return
}
