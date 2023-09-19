package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gmc/config"
	"gmc/filestore"
	"gmc/filestore/util"
)

func FileGet(exec string, cfg *config.Config, cmd string, args []string) {
	flagset := flag.NewFlagSet(cmd, flag.ExitOnError)
	out := flagset.String("out", "", "save location")
	flagset.SetOutput(os.Stdout)
	flagset.Usage = func() {
		fmt.Printf("Usage: %s %s <filename>\n",
			exec, cmd)
		flagset.PrintDefaults()
		os.Exit(1)
	}
	flagset.Parse(args)

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr,
			"A filename is required and a save destination is optional\n")
		flagset.Usage()
		os.Exit(1)
	}

	fs, err := filestore.New(cfg.FileStore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		os.Exit(1)
	}

	var file *util.File
	outpath := *out
	// if the outpath is empty, the file is saved in the cwd.
	if outpath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
		file, err = fs.GetFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
		if file == nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, errors.New("file not found"))
			os.Exit(1)
		}
		outpath = filepath.Join(cwd, (file).Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
	} else {
		file, err = fs.GetFile(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}
		outpath = filepath.Join(*out, (file).Name)
	}

	f, err := os.Create(outpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		os.Exit(1)
	}
	defer f.Close()
	_, err = io.Copy(f, file.Content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		os.Exit(1)
	}
}
