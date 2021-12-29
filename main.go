package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s: required argument missing\n", os.Args[0])
		printDefaultUsage()
		os.Exit(1)
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "start":
		startCommand()

	case "--help", "help":
		printDefaultUsage()
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a recognized command\n",
			os.Args[0], os.Args[1])
		printDefaultUsage()
		os.Exit(1)
	}
	os.Exit(0)
}

func printDefaultUsage() {
	fmt.Printf("Usage: %s <command> ...\n", os.Args[0])
	fmt.Printf("See '%s <command> --help' for information", os.Args[0])
	fmt.Printf(" on a specific command\n")
	fmt.Printf("valid commands:\n")
	fmt.Printf("    start    start HTTP server\n")
	fmt.Printf("    stop     stop running HTTP server\n")
}
