package cmd

import (
	"fmt"
	"os"
	"strings"

	"gmc/config"
	"gmc/db"
)

func KeywordCommand(cfg *config.Config, exec string, cmd string, args []string) int {
	printUsage := func() {
		fmt.Printf("Usage: %s [args] %s <subcommand> ...\n", exec, cmd)
		fmt.Printf("Subcommands:\n")
		fmt.Printf("  list, ls\n")
		fmt.Printf("      list keywords\n")
		fmt.Printf("  add <keywords ...>\n")
		fmt.Printf("      add new keyword(s)\n")
		fmt.Printf("  del <keywords ...>\n")
		fmt.Printf("      remove existing keyword(s)\n")
	}

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s %s: subcommand missing\n", exec, cmd)
		printUsage()
		return 1
	}

	subcmd := strings.ToLower(args[0])
	switch subcmd {
	case "-help", "--help", "help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "%s %s '%s' is not a recognized subcommand\n",
			exec, cmd, subcmd)
		printUsage()
		return 1

	case "list", "ls":
		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			return 1
		}

		kws, err := db.ListKeywords()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			return 1
		}

		for _, kw := range kws {
			fmt.Printf("%s\n", kw)
		}

	case "add":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s: new keyword name required\n", exec)
			return 1
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			return 1
		}

		if err := db.AddKeywords(args[1:]...); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			return 1
		}

	case "del", "rm", "delete", "remove":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s: keywords to remove required\n", exec)
			return 1
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			return 1
		}

		if err := db.DeleteKeywords(args[1:]...); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			return 1
		}
	}

	return 0
}
