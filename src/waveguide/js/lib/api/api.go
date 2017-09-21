package api

import (
	"github.com/gopherjs/gopherjs/js"
	"waveguide/lib/model"
)

type ApiProxy struct {
}

func (ap *ApiProxy) FindVideo() *js.Object {
	return js.MakeWrapper(&model.VideoInfo{})
}

func New() *js.Object {
	return js.MakeWrapper(&ApiProxy{})
}
