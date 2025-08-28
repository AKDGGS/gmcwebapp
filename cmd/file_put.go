package cmd

import (
	"flag"
	"fmt"
	"mime"
	"os"
	"path/filepath"

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
		fmt.Printf("Usage: %s %s <inventory flag> <inventory id> <filename ...>\n",
			exec, cmd)
		flagset.PrintDefaults()
	}
	if err := flagset.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	if *borehole_id == 0 && *inventory_id == 0 && *outcrop_id == 0 && *prospect_id == 0 && *well_id == 0 {
		fmt.Fprintf(os.Stderr, "no inventory flag (-borehole_id, -inventory_id, -outcrop_id, -prospect_id, -well_id) provided\n")
		flagset.Usage()
		return 1
	}

	if len(flagset.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "no filenames provided\n")
		flagset.Usage()
		return 1
	}

	db, err := db.New(cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		flagset.Usage()
		return 1
	}

	fs, err := filestore.NewFileStores(cfg.FileStores)
	if err != nil {
		return 1
	}

	exit_code := 0
	for _, filename := range flagset.Args() {
		file_info, err := os.Stat(filename)
		if err != nil || file_info.Size() == 0 {
			fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
			exit_code = 1
			continue
		}

		file := model.File{
			Name:         file_info.Name(),
			Size:         file_info.Size(),
			BoreholeIDs:  []int{*borehole_id},
			InventoryIDs: []int{*inventory_id},
			OutcropIDs:   []int{*outcrop_id},
			ProspectIDs:  []int{*prospect_id},
			WellIDs:      []int{*well_id},
		}

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
			fmt.Fprintf(os.Stderr, "%s\n", err)
			exit_code = 1
		}
	}
	return exit_code
}
