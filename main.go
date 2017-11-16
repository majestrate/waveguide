package main

import (
	"fmt"
	"os"
	"strings"
	"waveguide/cmd/cdn"
	"waveguide/cmd/frontend"
	"waveguide/cmd/worker"
	"waveguide/lib/version"
)

func printUsage() {
	fmt.Printf("usage: %s [frontend|worker|cdn]", os.Args[0])
	fmt.Println()
}

func main() {
	fmt.Printf("%s starting up", version.Version)
	fmt.Println()
	var mode string
	if len(os.Args) > 1 {
		mode = strings.ToUpper(os.Args[1])
	}
	switch mode {
	case "FRONTEND":
		frontend.Run()
	case "WORKER":
		worker.Run()
	case "CDN":
		cdn.Run()
	default:
		go frontend.Run()
		go worker.Run()
		cdn.Run()
	}
}
