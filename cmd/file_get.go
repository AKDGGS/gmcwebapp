package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gmc/config"
	"gmc/filestore"
	"gmc/filestore/util"
)

func FileGet(exec string, cfg *config.Config, cmd string, args []string) int {
	flagset := flag.NewFlagSet(cmd, flag.ContinueOnError)
	out := flagset.String("out", "", "save location")
	flagset.SetOutput(os.Stdout)
	flagset.Usage = func() {
		fmt.Printf("Usage: %s %s <filename>\n",
			exec, cmd)
		flagset.PrintDefaults()
	}
	if err := flagset.Parse(args); err == flag.ErrHelp {
		flagset.Usage()
		return 0
	}

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr,
			"a filename is required and a save destination is optional\n")
		flagset.Usage()
		return 1
	}

	if cfg.FileStore == nil {
		fmt.Fprintf(os.Stderr, "no file store configured\n")
		return 1
	}

	fs, err := filestore.New(*cfg.FileStore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	var file *util.File
	outpath := *out
	// if the outpath is empty, the file is saved in the cwd.
	if outpath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}
		file, err = fs.GetFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}
		if file == nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, fmt.Errorf("file not found"))
			return 1
		}
		outpath = filepath.Join(cwd, (file).Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}
	} else {
		file, err = fs.GetFile(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			return 1
		}
		outpath = filepath.Join(*out, (file).Name)
	}

	f, err := os.Create(outpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}
	defer f.Close()
	_, err = io.Copy(f, file.Content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	return 0
}
