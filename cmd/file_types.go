package cmd

import (
	"fmt"
	"os"
	"strings"

	"gmc/config"
	"gmc/db"
)

func FileTypesCommand(cfg *config.Config, exec string, cmd string, args []string) int {
	printUsage := func() {
		fmt.Printf("Usage: %s [args] %s <subcommand> ...\n", exec, cmd)
		fmt.Printf("Subcommands:\n")
		fmt.Printf("  list, ls\n")
		fmt.Printf("      list file types\n")
		fmt.Printf("  add <file types ...>\n")
		fmt.Printf("      add new file type(s)\n")
		fmt.Printf("  del <file types ...>\n")
		fmt.Printf("      remove existing file type(s)\n")
		fmt.Printf("  rename <oldname> <newname>\n")
		fmt.Printf("      rename an existing file type\n")
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
		db, err := db.New(cfg.Database)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}

		kws, err := db.ListFileTypes()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}

		for _, kw := range kws {
			fmt.Printf("%s\n", kw)
		}

	case "add":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s: new file type name required\n", exec)
			return 1
		}

		db, err := db.New(cfg.Database)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}

		ft, err := db.ListFileTypes()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}
		if len(ft)+(len(args)-1) > 128 {
			fmt.Fprintf(os.Stderr, "%s: cannot add more than 128 file types\n", exec)
			return 1
		}

		if err := db.AddFileTypes(args[1:]...); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}

	case "del", "rm", "delete", "remove":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s: file types to remove required\n", exec)
			return 1
		}

		db, err := db.New(cfg.Database)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}

		if err := db.DeleteFileTypes(args[1:]...); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}
	case "rename":
		if len(args) != 3 {
			fmt.Fprintf(os.Stderr, "%s: old and new file type names required\n", exec)
			return 1
		}

		db, err := db.New(cfg.Database)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}

		if err := db.RenameFileTypes(args[1], args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}
	}

	return 0
}
