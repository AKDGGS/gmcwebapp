package cmd

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"

	"gmc/config"
	"gmc/db"
	"gmc/db/model"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func TokenCommand(cfg *config.Config, exec string, cmd string, args []string) {
	printUsage := func() {
		fmt.Printf("Usage: %s [args] %s <subcommand> ...\n", exec, cmd)
		fmt.Printf("Subcommands:\n")
		fmt.Printf("  list, ls\n")
		fmt.Printf("      list personal access tokens\n")
		fmt.Printf("  delete, del, rm <id>\n")
		fmt.Printf("      remove a personal access token by ID\n")
		fmt.Printf("  create, add [args]\n")
		fmt.Printf("      create a personal access token\n")
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

		tokens, err := db.ListTokens()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

		for _, tk := range tokens {
			fmt.Printf("%d %s\n", tk.ID, tk.Description)
		}

	case "delete", "del", "rm":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s: token removal requires an ID\n", exec)
			os.Exit(1)
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: invalid token ID\n", exec)
			os.Exit(1)
		}

		err = db.DeleteToken(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

	case "create", "add":
		ffs := flag.NewFlagSet("token create", flag.ExitOnError)
		ffs.SetOutput(os.Stdout)
		ffs.Usage = func() {
			fmt.Fprintf(ffs.Output(), "Usage: %s %s %s [args]\n", exec, cmd, subcmd)
			fmt.Fprintf(ffs.Output(), "Arguments:\n")
			ffs.PrintDefaults()
			fmt.Fprintf(ffs.Output(), "\n")
			fmt.Fprintf(ffs.Output(), "If image is not specified, new token is ")
			fmt.Fprintf(ffs.Output(), "written as text to the console.\n")
		}
		img := ffs.String("image", "", "write new token as QR code image to path")
		width := ffs.Int("width", 750, "width of QR code image")
		height := ffs.Int("height", 750, "height of QR code image")
		desc := ffs.String("description", "", "description of new token (required)")
		ffs.Parse(args[1:])

		if desc == nil || *desc == "" {
			fmt.Fprintf(os.Stderr, "%s: token create requires description\n",
				os.Args[0])
			ffs.Usage()
			os.Exit(1)
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
			os.Exit(1)
		}

		tk := &model.Token{Description: *desc}
		tk.Generate()

		err = db.CreateToken(tk)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

		if img == nil || *img == "" {
			fmt.Printf("Generated token: %s\n", tk.Token)
			os.Exit(0)
		}

		qrcode, err := qr.Encode(tk.Token, qr.H, qr.AlphaNumeric)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
		qrcode, err = barcode.Scale(qrcode, *width, *height)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
		fout, err := os.Create(*img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
		defer fout.Close()
		err = png.Encode(fout, qrcode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
	}
}
