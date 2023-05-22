package cmd

import (
	"fmt"
	"os"
	"strings"

	"gmc/config"
)

func FileCommand(cfg *config.Config, exec string, cmd string, args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s %s: subcommand missing\n", exec, cmd)
		printFileUsage(exec, cmd)
		os.Exit(1)
	}

	switch subcmd := strings.ToLower(args[0]); subcmd {
	case "--help", "help":
		printFileUsage(exec, cmd)
	case "put":
		FilePut(exec, cfg, fmt.Sprintf("%s %s", cmd, subcmd), args[1:])
	case "get":
		FileGet(exec, cfg, fmt.Sprintf("%s %s", cmd, subcmd), args[1:])
	default:
		fmt.Fprintf(os.Stderr, "%s %s '%s' is not a recognized subcommand\n",
			exec, cmd, subcmd)
		printFileUsage(exec, cmd)
		os.Exit(1)
	}
	return nil
}

func printFileUsage(exec, cmd string) {
	fmt.Printf("Usage: %s %s <subcommand> ...\n", exec, cmd)
	fmt.Printf("Subcommands:\n")
	fmt.Printf("  put [args] <filename ...>\n")
	fmt.Printf("      upload file to filestore\n")
	fmt.Printf("  get [args] <filename ...>\n")
	fmt.Printf("      download file from filestore\n")
}
