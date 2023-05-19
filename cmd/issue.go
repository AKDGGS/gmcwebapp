package cmd

import (
	"fmt"
	"os"
	"strings"

	"gmc/config"
	"gmc/db"
)

func IssueCommand(cfg *config.Config, exec string, cmd string, args []string) {
	printUsage := func() {
		fmt.Printf("Usage: %s [args] %s <subcommand> ...\n", exec, cmd)
		fmt.Printf("Subcommands:\n")
		fmt.Printf("  list, ls\n")
		fmt.Printf("      list issues\n")
		fmt.Printf("  add <issues ...>\n")
		fmt.Printf("      add new quality issue(s)\n")
		fmt.Printf("  del <issues ...>\n")
		fmt.Printf("      remove quality issue(s)\n")
	}

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s %s: subcommand missing\n", exec, cmd)
		printUsage()
		os.Exit(1)
	}

	subcmd := strings.ToLower(args[0])
	switch subcmd {
	case "--help", "help":
		printUsage()
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "%s %s '%s' is not a recognized subcommand\n",
			exec, cmd, subcmd)
		printUsage()
		os.Exit(1)

	case "list", "ls":
		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

		iss, err := db.ListIssues()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

		for _, is := range iss {
			fmt.Printf("%s\n", is)
		}

	case "add":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s: new issue name required\n", exec)
			os.Exit(1)
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

		if err := db.AddIssues(args[1:]...); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
	case "del", "rm", "delete", "remove":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s: issues to remove required\n", exec)
			os.Exit(1)
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

		if err := db.DeleteIssues(args[1:]...); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
	}
}
