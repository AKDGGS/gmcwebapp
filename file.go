package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"math/big"
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
		fmt.Printf("  list, ls\n")
		fmt.Printf("      list issues\n")
		fmt.Printf("  put <filename ...>\n")
		fmt.Printf("      add new quality issue(s)\n")
		fmt.Printf("  get <filename ...>\n")
		fmt.Printf("      remove quality issue(s)\n")
	}
	fileFlagSet := flag.NewFlagSet("file", flag.ContinueOnError)
	fileFlagSet.Usage = func() {
		printUsage()
	}
	fileFlagSet.Parse(args)

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s %s: subcommand missing\n", exec, cmd)
		printUsage()
		os.Exit(1)
	}

	subcmd := strings.ToLower(fileFlagSet.Arg(0))
	var subcmdArgs []string
	if len(fileFlagSet.Args()) > 1 {
		subcmdArgs = fileFlagSet.Args()[1:]
	}

	switch subcmd {
	default:
		flagset := flag.NewFlagSet(cmd, flag.ExitOnError)
		flagset.SetOutput(os.Stdout)
		flagset.Usage = func() {
			fmt.Printf("Usage: %s %s %s <filename>\n",
				exec, cmd, subcmd)
			flagset.PrintDefaults()
		}
	case "put":
		flagset := flag.NewFlagSet(cmd, flag.ExitOnError)
		well_id := flagset.Int("well_id", 0, "a well ID")
		flagset.SetOutput(os.Stdout)
		flagset.Usage = func() {
			fmt.Printf("Usage: %s %s %s <filename>\n",
				exec, cmd, subcmd)
			flagset.PrintDefaults()
		}
		flagset.Parse(subcmdArgs)

		well := *well_id

		if well == 0 {
			subcmdArgs = subcmdArgs[0:]
		} else {
			subcmdArgs = subcmdArgs[2:]
		}

		var fileArr []model.File
		for _, filename := range subcmdArgs {
			fileInfo, err := os.Stat(filename)
			if err != nil || fileInfo.Size() == 0 {
				fmt.Println("At least one of the files doesn't exist or has a size of zero.")
				break
			}

			// temporary code until we decide what to do with the MD5.
			rand.Seed(time.Now().UnixNano())
			MD5String := strconv.FormatInt(rand.Int63(), 10)

			file := model.File{}
			file.Name = fileInfo.Name()
			file.Size = fileInfo.Size()
			file.MD5 = MD5String
			fileArr = append(fileArr, file)
		}

		db, err := db.New(cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
			os.Exit(1)
		}

		for i, _ := range fileArr {
			err = db.PutFile(&fileArr[i], func() error {
				fs, err := filestore.New(cfg.FileStore)
				if err != nil {
					return nil
				}

				for _, filename := range subcmdArgs {
					file, err := os.Open(filename)
					if err != nil {
						return err
					}

					stat, err := file.Stat()
					if err != nil {
						return err
					}

					mt := mime.TypeByExtension(filepath.Ext(stat.Name()))
					if mt == "" {
						mt = "application/octet-stream"
					}

					md_b := big.NewInt(stat.ModTime().UnixMicro()).Bytes()
					sz_b := big.NewInt(stat.Size()).Bytes()

					f := &fsutil.File{
						Name:         stat.Name(),
						Size:         stat.Size(),
						LastModified: stat.ModTime(),
						ContentType:  mt,
						Content:      file,
						ETag: fmt.Sprintf("%s-%s",
							base64.RawStdEncoding.EncodeToString(md_b),
							base64.RawStdEncoding.EncodeToString(sz_b),
						),
					}
					fs.PutFile(f)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		// fs, err := filestore.New(cfg.FileStore)
		// if err != nil {
		// 	return nil
		// }
		//
		// for _, filename := range subcmdArgs {
		// 	file, err := os.Open(filename)
		// 	if err != nil {
		// 		return err
		// 	}
		//
		// 	stat, err := file.Stat()
		// 	if err != nil {
		// 		return err
		// 	}
		//
		// 	mt := mime.TypeByExtension(filepath.Ext(stat.Name()))
		// 	if mt == "" {
		// 		mt = "application/octet-stream"
		// 	}
		//
		// 	md_b := big.NewInt(stat.ModTime().UnixMicro()).Bytes()
		// 	sz_b := big.NewInt(stat.Size()).Bytes()
		//
		// 	f := &fsutil.File{
		// 		Name:         stat.Name(),
		// 		Size:         stat.Size(),
		// 		LastModified: stat.ModTime(),
		// 		ContentType:  mt,
		// 		Content:      file,
		// 		ETag: fmt.Sprintf("%s-%s",
		// 			base64.RawStdEncoding.EncodeToString(md_b),
		// 			base64.RawStdEncoding.EncodeToString(sz_b),
		// 		),
		// 	}
		// 	fs.PutFile(f)
		// }

	case "get":
		flagset := flag.NewFlagSet(cmd, flag.ExitOnError)
		flagset.SetOutput(os.Stdout)
		out := flagset.String("out", "", "output file")
		flagset.Usage = func() {
			fmt.Printf("Usage: %s %s %s [args] <filename>\n",
				exec, cmd, subcmd)
			flagset.PrintDefaults()
		}
		flagset.Parse(subcmdArgs)

		if len(subcmdArgs) < 1 || len(args)%2 != 0 {
			fmt.Fprintf(os.Stderr,
				"A subcommand and filename are required\n")
			flagset.Usage()
			os.Exit(1)
		}

		fs, err := filestore.New(cfg.FileStore)
		if err != nil {
			return nil
		}

		var file *util.File
		outputPath := *out
		if outputPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			outputPath = filepath.Join(cwd, filepath.Base(subcmdArgs[0]))
			file, err = fs.GetFile(subcmdArgs[0])
			if err != nil {
				return err
			}
		} else {
			outputPath = filepath.Join(*out, filepath.Base(subcmdArgs[2]))
			file, err = fs.GetFile(fileFlagSet.Arg(3))
			if err != nil {
				return err
			}
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, file.Content)
		if err != nil {
			return err
		}

	}
	return nil
}
