package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gmc/assets"
	"gmc/config"
)

func main() {
	exec := os.Args[0]
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = func() {
		fmt.Printf("Usage: %s [args] <command> ...\n", exec)
		fmt.Printf("Arguments:\n")
		flag.PrintDefaults()
		fmt.Printf("Commands:\n")
		fmt.Printf("  start, server\n")
		fmt.Printf("      start HTTP server\n")
		fmt.Printf("  db, database\n")
		fmt.Printf("      initialize, verify, or drop a database\n")
		fmt.Printf("  genkey, keygen\n")
		fmt.Printf("      generate a random session key\n")
		fmt.Printf("  token, tk\n")
		fmt.Printf("      manage personal access tokens\n")
		fmt.Printf("  keyword, kw\n")
		fmt.Printf("      manage keywords\n")

		fmt.Printf("See '%s <command> --help' for information ", exec)
		fmt.Printf("on a specific command\n")
	}
	cpath := flag.String("conf", "", "path to configuration")
	assetpath := flag.String("assets", "", "override embedded assets with assets from path")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "%s: required argument missing\n", exec)
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.Load(*cpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		os.Exit(1)
	}

	if assetpath != nil {
		assets.SetExternal(*assetpath)
	}

	cmd := strings.ToLower(flag.Arg(0))
	switch cmd {
	case "database", "db":
		databaseCommand(cfg, exec, cmd, flag.Args()[1:])

	case "server", "start":
		serverCommand(cfg, exec)

	case "genkey", "keygen":
		genkeyCommand()

	case "token", "tk":
		tokenCommand(cfg, exec, cmd, flag.Args()[1:])

	case "keyword", "kw":
		keywordCommand(cfg, exec, cmd, flag.Args()[1:])

	case "--help", "help":
		flag.Usage()
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a recognized command\n", exec, cmd)
		flag.Usage()
		os.Exit(1)
	}
	os.Exit(0)
}
