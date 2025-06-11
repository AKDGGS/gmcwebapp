package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"gmc/config"
	"gmc/db"
	"gmc/db/model"
	"gmc/filestore"
)

func FileGet(exec string, cfg *config.Config, cmd string, args []string) int {
	flagset := flag.NewFlagSet(cmd, flag.ContinueOnError)
	out := flagset.String("out", "", "output directory")
	flagset.SetOutput(os.Stdout)
	flagset.Usage = func() {
		fmt.Printf("Usage: %s %s -out <output directory> <file ids ...>\n",
			exec, cmd)
		flagset.PrintDefaults()
	}
	if err := flagset.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	if *out == "" {
		fmt.Fprintln(os.Stderr, "no output directory provided")
		flagset.Usage()
		return 1
	}

	fi, err := os.Stat(*out)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: output directory does not exist\n", exec)
		} else if os.IsPermission(err) {
			fmt.Fprintf(os.Stderr, "%s: permission denied accessing output directory\n", exec)
		} else {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		}
		return 1
	}
	if !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "%s: %s is not a directory\n", exec, *out)
		return 1
	}

	if flagset.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "no file ids provided")
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
	db, err := db.New(cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		return 1
	}

	var valid_db_files []*model.File

	for _, fid := range flagset.Args() {
		id, err := strconv.Atoi(fid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %d %s\n", exec, id, err)
			return 1
		}
		db_file, err := db.GetFile(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %d %s\n", exec, id, err)
			return 1
		}
		if db_file == nil {
			fmt.Fprintf(os.Stderr, "%s: %d file does not exist\n", exec, id)
			return 1
		}
		valid_db_files = append(valid_db_files, db_file)
	}
	exit_code := 0
	for _, db_file := range valid_db_files {
		file, err := fs.GetFile(fmt.Sprintf("%d/%s", db_file.ID, db_file.Name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			exit_code = 1
			continue
		}
		f, err := os.OpenFile(filepath.Join(*out, db_file.Name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			if os.IsExist(err) {
				fmt.Fprintf(os.Stderr, "%s: a file with the same name already exists: %s\n", exec, db_file.Name)
			} else {
				fmt.Fprintf(os.Stderr, "%s: error creating file: %s\n", exec, err)
			}
			exit_code = 1
			continue
		}
		defer f.Close()
		_, err = io.Copy(f, file.Content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s\n", exec, err)
			exit_code = 1
			continue
		}
	}
	return exit_code
}
