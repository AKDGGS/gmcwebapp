package main

import (
	"flag"
	"fmt"
	"gmc/config"
	"gmc/db"
	"os"
	"strings"
)

func databaseCommand() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s database: subcommand missing\n", os.Args[0])
		printDatabaseUsage()
		os.Exit(1)
	}

	subcmd := strings.ToLower(os.Args[2])
	switch subcmd {
	case "--help", "help":
		printDatabaseUsage()
		os.Exit(0)

	case "verify":
		cmd := flag.NewFlagSet("database verify", flag.ExitOnError)
		cmd.SetOutput(os.Stdout)
		cmd.Usage = func() {
			fmt.Printf("Verifies a database connection with a simple query.\n\n")
			fmt.Printf("Usage: %s database %s [args]\n", os.Args[0], subcmd)
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

		if err := db.Verify(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

	case "drop":
		cmd := flag.NewFlagSet("database drop", flag.ExitOnError)
		cmd.SetOutput(os.Stdout)
		cmd.Usage = func() {
			fmt.Printf("Drops entire database. Don't use this, seriously.\n\n")
			fmt.Printf("Usage: %s database %s [args]\n", os.Args[0], subcmd)
			cmd.PrintDefaults()
		}
		cpath := cmd.String("conf", "", "path to configuration")
		sure := cmd.Bool(
			"yesreally", false,
			"Are you really sure you want to do this?",
		)
		cmd.Parse(os.Args[3:])

		if !*sure {
			fmt.Fprintf(os.Stderr, "%s: refusing to drop since ", os.Args[0])
			fmt.Fprintf(os.Stderr, "you're not really sure\n")
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

		if err := db.Drop(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}


	case "init", "initialize":
		cmd := flag.NewFlagSet("database init", flag.ExitOnError)
		cmd.SetOutput(os.Stdout)
		cmd.Usage = func() {
			fmt.Printf("Initializes an empty database for use with the GMC ")
			fmt.Printf("application. Required when\nstarting a new database.\n\n")
			fmt.Printf("Usage: %s database %s [args]\n", os.Args[0], subcmd)
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

		if err := db.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}
	}
}

func printDatabaseUsage() {
	fmt.Printf("Usage: %s database <subcommand> ...\n", os.Args[0])
	fmt.Printf("See '%s database <subcommand> --help' for", os.Args[0])
	fmt.Printf(" information on a specific command\n")
	fmt.Printf("valid commands:\n")
	fmt.Printf("    init    initialize an empty database\n")
	fmt.Printf("    verify  verify database configuration\n")
	fmt.Printf("    drop    drops entire database (DANGER!)\n")

}
