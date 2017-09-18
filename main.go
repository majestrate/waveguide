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
	fmt.Printf("usage: %s [worker|frontend]", os.Args[0])
	fmt.Println()
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		return
	}
	args := os.Args[1:]
	mode := strings.ToUpper(args[0])
	fmt.Printf("%s starting up %s", version.Version, mode)
	fmt.Println()
	switch mode {
	case "WORKER":
		worker.Run()
	case "FRONTEND":
		frontend.Run()
	default:
		printUsage()
	}
}
