package main

import (
	"github.com/gopherjs/gopherjs/js"
	"waveguide/js/lib/api"
)

func main() {
	js.Global.Set("waveguide", map[string]interface{}{
		"New": api.New,
	})
}
