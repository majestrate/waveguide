package api

import (
	"io"
	"mime/multipart"
	"net/http"
)

type Request struct {
	c   io.Closer
	mpr *multipart.Reader
}

func (r *Request) WalkFiles(f func(string, io.Reader) error) (err error) {
	var part *multipart.Part
	for err == nil {
		part, err = r.mpr.NextPart()
		if err == io.EOF {
			err = nil
			break
		}
		err = f(part.FileName(), part)
		part.Close()
	}
	return
}

/** impelements io.Closer */
func (r *Request) Close() error {
	return r.c.Close()
}

/** TransformRequest transforms an http.Request to a pomf request */
func TransformRequest(r *http.Request) (req *Request, err error) {
	req = &Request{
		c: r.Body,
	}
	req.mpr, err = r.MultipartReader()
	if err == nil {
		return
	}
	req = nil
	return
}
