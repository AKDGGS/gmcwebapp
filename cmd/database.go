package cmd

import (
	"fmt"
	"os"
	"strings"

	"gmc/config"
	"gmc/db"
)

func DatabaseCommand(cfg *config.Config, exec string, cmd string, args []string) {
	printUsage := func() {
		fmt.Printf("Usage: %s [args] %s <subcommand> ...\n", exec, cmd)
		fmt.Printf("Subcommands:\n")
		fmt.Printf("  initialize, init\n")
		fmt.Printf("      initialize an empty database with a fresh schema\n")
		fmt.Printf("  verify\n")
		fmt.Printf("      verify database configuration\n")
		fmt.Printf("  drop\n")
		fmt.Printf("      drops entire database schema (DANGER!)\n")
	}

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s %s: subcommand missing\n", exec, cmd)
		printUsage()
		os.Exit(1)
	}

	subcmd := strings.ToLower(args[0])
	switch subcmd {
	case "initialize", "init":
		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		if err := db.SchemaInit(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

	case "verify":
		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		if err := db.Verify(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		fmt.Printf("Verification successful.\n")

	case "drop":
		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		var confirm string
		fmt.Printf("Doing this will empty the current database, ")
		fmt.Printf("removing the existing schema.\n")
		fmt.Printf("Are you sure you want to proceed (yes/no)? ")
		fmt.Scanf("%s", &confirm)

		if confirm == "yes" {
			if err := db.SchemaDrop(); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Printf("aborted\n")
		}

	case "--help", "help":
		printUsage()
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "%s %s '%s' is not a recognized subcommand\n",
			exec, cmd, subcmd)
		printUsage()
		os.Exit(1)
	}
}
