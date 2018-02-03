package api

import (
	"fmt"
)

type Request struct {
	Method string
	Args   map[string]interface{}
}

func (r *Request) Get(key string, fallback interface{}) interface{} {
	v, ok := r.Args[key]
	if ok {
		return v
	}
	return fallback
}

func (r *Request) GetInt(key string, fallback int64) int64 {
	f, ok := r.Args[key]
	if ok {
		v, ok := f.(float64)
		if ok {
			return int64(v)
		}
	}
	return fallback
}

func (r *Request) GetString(key, fallback string) string {
	v, ok := r.Args[key]
	if ok {
		return fmt.Sprintf("%s", v)
	}
	return fallback
}
