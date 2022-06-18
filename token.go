package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gmc/config"
	"gmc/db"
)

func tokenCommand(rootcmd string) {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s %s: subcommand missing\n", os.Args[0], rootcmd)
		printTokenUsage(rootcmd)
		os.Exit(1)
	}

	subcmd := strings.ToLower(os.Args[2])
	switch subcmd {
	case "--help", "help":
		printTokenUsage(rootcmd)
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "%s %s '%s' is not a recognized subcommand\n",
			os.Args[0], rootcmd, subcmd)
		printTokenUsage(rootcmd)
		os.Exit(1)

	case "list", "ls":
		cmd := flag.NewFlagSet("token list", flag.ExitOnError)
		cmd.SetOutput(os.Stdout)
		cmd.Usage = func() {
			fmt.Printf("Lists all available tokens\n\n")
			fmt.Printf("Usage: %s %s %s [args]\n", os.Args[0], rootcmd, subcmd)
			cmd.PrintDefaults()
		}
		cpath := cmd.String("conf", "", "path to configuration")
		cmd.Parse(os.Args[3:])

		cfg, err := config.Load(*cpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		tokens, err := db.ListTokens()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		for _, tk := range tokens {
			fmt.Printf("%d %s\n", tk.ID, tk.Description)
		}

	case "delete", "del", "rm":
		cmd := flag.NewFlagSet("token delete", flag.ExitOnError)
		cmd.SetOutput(os.Stdout)
		cmd.Usage = func() {
			fmt.Printf("Deletes a token by ID\n\n")
			fmt.Printf("Usage: %s %s %s [args] [ID]\n", os.Args[0], rootcmd, subcmd)
			cmd.PrintDefaults()
		}
		cpath := cmd.String("conf", "", "path to configuration")
		cmd.Parse(os.Args[3:])

		if cmd.NArg() < 1 {
			fmt.Fprintf(os.Stderr, "%s: token delete requires an ID\n", os.Args[0])
			os.Exit(1)
		}

		cfg, err := config.Load(*cpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		id, err := strconv.Atoi(cmd.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: invalid token ID\n", os.Args[0])
			os.Exit(1)
		}

		err = db.DeleteToken(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}
	}
}

func printTokenUsage(rootcmd string) {
	fmt.Printf("Usage: %s %s <subcommand> ...\n", os.Args[0], rootcmd)
	fmt.Printf("See '%s %s <subcommand> --help' for", os.Args[0], rootcmd)
	fmt.Printf(" information on a specific command\n")
	fmt.Printf("valid commands:\n")
	fmt.Printf("    list, ls         list personal access tokens\n")
	fmt.Printf("    delete, del, rm  remove a personal access token by ID\n")
}
