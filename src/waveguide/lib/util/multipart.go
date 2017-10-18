package util

import (
	"io"
	"mime/multipart"
)

type MimePart struct {
	Body     io.ReadCloser
	PartName string
}

type MultipartPipe struct {
	Parts []MimePart
	pr    *io.PipeReader
	pw    *io.PipeWriter
}

func (p *MultipartPipe) Read(data []byte) (int, error) {
	return p.pr.Read(data)
}

func (p *MultipartPipe) Close() error {
	return p.pr.Close()
}

func (p *MultipartPipe) Run() {
	var buff [65536]byte
	mpw := multipart.NewWriter(p.pw)
	for _, info := range p.Parts {
		body, err := mpw.CreateFormFile(info.PartName, "")
		if err == nil {
			io.CopyBuffer(body, info.Body, buff[:])
			info.Body.Close()
		}
	}
	mpw.Close()
	p.pw.Close()
}

func NewMultipartPipe(parts []MimePart) *MultipartPipe {
	pr, pw := io.Pipe()
	p := &MultipartPipe{
		Parts: parts,
		pr:    pr,
		pw:    pw,
	}
	go p.Run()
	return p
}
