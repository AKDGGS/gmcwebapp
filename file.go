package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gmc/config"
	"gmc/db"
	"gmc/db/model"
	"gmc/filestore"
	"gmc/filestore/util"
	fsutil "gmc/filestore/util"
)

func fileCommand(cfg *config.Config, exec string, cmd string, args []string) error {
	printUsage := func() {
		fmt.Printf("Usage: %s %s <subcommand> ...\n", exec, cmd)
		fmt.Printf("Subcommands:\n")
		fmt.Printf("  put [args] <filename ...>\n")
		fmt.Printf("      upload file to filestore\n")
		fmt.Printf("  get [args] <filename ...>\n")
		fmt.Printf("      download file from filestore\n")
	}

	file_flagset := flag.NewFlagSet("file", flag.ContinueOnError)
	file_flagset.SetOutput(ioutil.Discard)
	err := file_flagset.Parse(args)
	if err != nil {
		printUsage()
		os.Exit(1)
	}
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s %s: subcommand missing\n", exec, cmd)
		printUsage()
		os.Exit(1)
	}

	subcmd := strings.ToLower(file_flagset.Arg(0))
	var subcmd_args []string
	if len(file_flagset.Args()) > 1 {
		subcmd_args = file_flagset.Args()[1:]
	}

	switch subcmd {
	case "--help", "help":
		printUsage()
		os.Exit(0)
	case "put":
		flagset := flag.NewFlagSet(cmd, flag.ExitOnError)
		well_id := flagset.Int("well_id", 0, "a well ID")
		flagset.SetOutput(os.Stdout)
		flagset.Usage = func() {
			fmt.Printf("Usage: %s %s %s <filename>\n",
				exec, cmd, subcmd)
			flagset.PrintDefaults()
			os.Exit(1)
		}
		flagset.Parse(subcmd_args)

		if *well_id == 0 {
			subcmd_args = subcmd_args[0:]
		} else {
			subcmd_args = subcmd_args[2:]
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

		fs, err := filestore.New(cfg.FileStore)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			return nil
		}

		for _, filename := range subcmd_args {
			file_info, err := os.Stat(filename)
			if err != nil || file_info.Size() == 0 {
				fmt.Println(filename, "does not exist or has a size of zero")
			} else {
				// temporary code until we decide what to do with the MD5.
				rand.Seed(time.Now().UnixNano())
				MD5 := strconv.FormatInt(rand.Int63(), 10)

				file := model.File{
					Name: file_info.Name(),
					Size: file_info.Size(),
					MD5:  MD5,
				}

				file.WellIDs = append(file.WellIDs, *well_id)

				// Add the file to the database
				err = db.PutFile(&file, func() error {
					file_obj, err := os.Open(file.Name)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
						return err
					}

					mt := mime.TypeByExtension(filepath.Ext(file_info.Name()))
					if mt == "" {
						mt = "application/octet-stream"
					}

					//Add the file to the filestore
					err = fs.PutFile(&fsutil.File{
						Name:         fmt.Sprintf("%d/%s", file.ID, file_info.Name()),
						Size:         file_info.Size(),
						LastModified: file_info.ModTime(),
						ContentType:  mt,
						Content:      file_obj,
					})
					if err != nil {
						return fmt.Errorf("error putting file in filestore: %w", err)
					}
					return nil
				})
				if err != nil {
					return fmt.Errorf("error putting file in database or the filestore: %w", err)
				}
			}
		}

	case "get":
		flagset := flag.NewFlagSet(cmd, flag.ExitOnError)
		flagset.SetOutput(os.Stdout)
		out := flagset.String("out", "", "output file")
		flagset.Usage = func() {
			fmt.Printf("Usage: %s %s %s [args] <filename>\n",
				exec, cmd, subcmd)
			flagset.PrintDefaults()
			os.Exit(1)
		}
		flagset.Parse(subcmd_args)

		if len(subcmd_args) < 1 || len(args)%2 != 0 {
			fmt.Fprintf(os.Stderr,
				"A filename is required and a save destination is optional\n")
			flagset.Usage()
			os.Exit(1)
		}

		fs, err := filestore.New(cfg.FileStore)
		if err != nil {
			return nil
		}

		var file *util.File
		outpath := *out
		if outpath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
				return err
			}
			outpath = filepath.Join(cwd, filepath.Base(subcmd_args[0]))
			file, err = fs.GetFile(subcmd_args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
				return err
			}
		} else {
			outpath = filepath.Join(*out, filepath.Base(subcmd_args[2]))
			file, err = fs.GetFile(file_flagset.Arg(3))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
				return err
			}
		}

		f, err := os.Create(outpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, file.Content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			return err
		}

	default:
		fmt.Fprintf(os.Stderr, "%s %s '%s' is not a recognized subcommand\n",
			exec, cmd, subcmd)
		printUsage()
		os.Exit(1)
	}
	return nil
}
