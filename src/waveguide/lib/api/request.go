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

func (r *Request) GetString(key, fallback string) string {
	v, ok := r.Args[key]
	if ok {
		return fmt.Sprintf("%s", v)
	}
	return fallback
}
