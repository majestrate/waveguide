package main

import (
	"fmt"
	"os"
	"strings"
	"waveguide/cmd/frontend"
	"waveguide/cmd/worker"
	"waveguide/lib/version"
)

func printUsage() {
	fmt.Printf("usage: %s [frontend|worker]", os.Args[0])
	fmt.Println()
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		return
	}
	fmt.Printf("%s starting up", version.Version)
	fmt.Println()
	mode := strings.ToUpper(os.Args[1])
	switch mode {
	case "FRONTEND":
		frontend.Run()
	case "WORKER":
		worker.Run()
	default:
		printUsage()
	}
}
