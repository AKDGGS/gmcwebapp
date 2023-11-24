package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gmc/assets"
	"gmc/cmd"
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
		fmt.Printf("  keyword, keywords, kw\n")
		fmt.Printf("      manage keywords\n")
		fmt.Printf("  issue, issues, iss\n")
		fmt.Printf("      manage quality issues\n")
		fmt.Printf("  file, f\n")
		fmt.Printf("      manage files\n")
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

	if *assetpath != "" {
		assets.SetExternal(*assetpath)
	}

	cm := strings.ToLower(flag.Arg(0))
	switch cm {
	case "database", "db":
		cmd.DatabaseCommand(cfg, exec, cm, flag.Args()[1:])

	case "server", "start":
		cmd.ServerCommand(cfg, exec)

	case "genkey", "keygen":
		cmd.GenKeyCommand()

	case "token", "tk":
		cmd.TokenCommand(cfg, exec, cm, flag.Args()[1:])

	case "keywords", "keyword", "kw":
		cmd.KeywordCommand(cfg, exec, cm, flag.Args()[1:])

	case "issues", "issue", "iss":
		cmd.IssueCommand(cfg, exec, cm, flag.Args()[1:])

	case "file", "f":
		cmd.FileCommand(cfg, exec, cm, flag.Args()[1:])

	case "--help", "help":
		flag.Usage()

	default:
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a recognized command\n", exec, cm)
		flag.Usage()
		os.Exit(1)
	}
	os.Exit(0)
}
