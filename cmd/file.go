package cmd

import (
	"fmt"
	"os"
	"strings"

	"gmc/config"
)

func FileCommand(cfg *config.Config, exec, cmd string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s %s: subcommand missing\n", exec, cmd)
		printFileUsage(exec, cmd)
		return 1
	}

	switch subcmd := strings.ToLower(args[0]); subcmd {
	case "-help", "--help", "help":
		printFileUsage(exec, cmd)
	case "put":
		return FilePut(exec, cfg, fmt.Sprintf("%s %s", cmd, subcmd), args[1:])
	case "get":
		return FileGet(exec, cfg, fmt.Sprintf("%s %s", cmd, subcmd), args[1:])
	default:
		fmt.Fprintf(os.Stderr, "%s %s '%s' is not a recognized subcommand\n",
			exec, cmd, subcmd)
		printFileUsage(exec, cmd)
		return 1
	}
	return 0
}

func printFileUsage(exec, cmd string) {
	fmt.Printf("Usage: %s %s <subcommand> ...\n", exec, cmd)
	fmt.Printf("Subcommands:\n")
	fmt.Printf("  put [args] <filename ...>\n")
	fmt.Printf("      upload files to filestore\n")
	fmt.Printf("  get [args] -out <output directory> <file ids ...>\n")
	fmt.Printf("      download file from filestore\n")
}
