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

func FilePut(exec string, cfg *config.Config, cmd string, args []string) int {
	flagset := flag.NewFlagSet(cmd, flag.ContinueOnError)
	borehole_id := flagset.Int("borehole_id", 0, "Borehole ID linked to file")
	inventory_id := flagset.Int("inventory_id", 0, "Inventory ID linked to file")
	outcrop_id := flagset.Int("outcrop_id", 0, "Outcrop ID linked to file")
	prospect_id := flagset.Int("prospect_id", 0, "Prospect ID linked to file")
	well_id := flagset.Int("well_id", 0, "Well ID linked to file")
	flagset.SetOutput(os.Stdout)
	flagset.Usage = func() {
		fmt.Printf("Usage: %s %s <filename>\n",
			exec, cmd)
		flagset.PrintDefaults()
	}
	if err := flagset.Parse(args); err == flag.ErrHelp {
		flagset.Usage()
		return 1
	}
	if *borehole_id == 0 && *inventory_id == 0 && *outcrop_id == 0 && *prospect_id == 0 && *well_id == 0 {
		fmt.Fprintf(os.Stderr, "linking flag (-well_id, -outcrop_id) is required\n")
		flagset.Usage()
		return 1
	}

	if len(flagset.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "filename required\n")
		flagset.Usage()
		return 1
	}

	filenames := flagset.Args()
	db, err := db.New(cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		flagset.Usage()
		return 1
	}

	fs, err := filestore.New(cfg.FileStore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		return 1
	}

	for _, filename := range filenames {
		file_info, err := os.Stat(filename)
		if err != nil || file_info.Size() == 0 {
			return 1
		}
		// temporary code until we decide what to do with the MD5.
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		MD5 := strconv.FormatInt(random.Int63(), 10)

		file := model.File{
			Name: file_info.Name(),
			Size: file_info.Size(),
			MD5:  MD5,
		}

		file.BoreholeIDs = append(file.BoreholeIDs, *borehole_id)
		file.InventoryIDs = append(file.InventoryIDs, *inventory_id)
		file.OutcropIDs = append(file.OutcropIDs, *outcrop_id)
		file.ProspectIDs = append(file.ProspectIDs, *prospect_id)
		file.WellIDs = append(file.WellIDs, *well_id)
		// Add the file to the database
		err = db.PutFile(&file, func() error {
			file_obj, err := os.Open(filename)
			if err != nil {
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
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return 1
		}
	}
	return 0
}
