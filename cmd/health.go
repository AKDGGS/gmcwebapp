package cmd

import (
	"fmt"
	"os"
	"strings"

	"gmc/config"
)

func HealthCommand(cfg *config.Config, exec string, cmd string, args []string) int {
	printUsage := func() {
		fmt.Printf("Usage: %s [args] %s <subcommand> ...\n", exec, cmd)
		fmt.Printf("Subcommands:\n")
		fmt.Printf("  url\n")
		fmt.Printf("      check all urls\n")
	}

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s %s: subcommand missing\n", exec, cmd)
		printUsage()
		return 1
	}

	subcmd := strings.ToLower(args[0])
	switch subcmd {
	case "url":
		return HealthURL(exec, cfg)

	case "-help", "--help", "help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "%s %s '%s' is not a recognized subcommand\n",
			exec, cmd, subcmd)
		printUsage()
		return 1
	}

	return 0
}
