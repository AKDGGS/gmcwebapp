package cmd

import (
	"flag"
	"fmt"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gmc/config"
	"gmc/db"
	"gmc/db/model"
	"gmc/filestore"
	fsutil "gmc/filestore/util"
)

func FilePut(exec string, cfg *config.Config, cmd string, args []string) {
	flagset := flag.NewFlagSet(cmd, flag.ExitOnError)
	well_id := flagset.Int("well_id", 0, "Well ID linked to file")
	flagset.SetOutput(os.Stdout)
	flagset.Usage = func() {
		fmt.Printf("Usage: %s %s <filename>\n",
			exec, cmd)
		flagset.PrintDefaults()
	}
	flagset.Parse(args)
	if *well_id == 0 {
		fmt.Fprintf(os.Stderr, "-well_id flag is required\n")
		flagset.Usage()
		os.Exit(1)
	}

	if len(flagset.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "filename required\n")
		flagset.Usage()
		os.Exit(1)
	}

	filenames := flagset.Args()
	db, err := db.New(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		flagset.Usage()
		os.Exit(1)
	}

	fs, err := filestore.New(cfg.FileStore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		os.Exit(1)
	}

	for _, filename := range filenames {
		file_info, err := os.Stat(filename)
		if err != nil || file_info.Size() == 0 {
			continue
		}
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
				os.Exit(1)
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
				return fmt.Errorf("Error putting file in filestore: %w", err)
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	}
}
