package main

import (
	"fmt"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"gmc/config"
	"gmc/db"
	"gmc/db/model"
	"gmc/filestore"
	fsutil "gmc/filestore/util"
)

type params struct {
	P_type        string
	Content_bytes int64
	ID            int32
}

func TestFile(t *testing.T) {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		t.Fatal("Error:", err)
	}
	cfg, err := config.Load(filepath.Join(config_dir, "gmc.test.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	db, err := db.New(cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	fs, err := filestore.New(cfg.FileStore)
	if err != nil {
		t.Fatal(err)
	}
	params := []params{
		{P_type: "borehole", Content_bytes: 0},
		{P_type: "inventory", Content_bytes: 100},
		{P_type: "outcrop", Content_bytes: 1000},
		{P_type: "prospect", Content_bytes: 10000},
		{P_type: "well", Content_bytes: 100000},
		{P_type: "", Content_bytes: 100},
	}
	for _, p := range params {
		f, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatal(err)
		}
		file_name := filepath.Base(f.Name())
		defer os.Remove(f.Name())
		defer f.Close()
		f_content_bytes := p.Content_bytes
		f_contents := make([]byte, f_content_bytes)
		for j := range f_contents {
			f_contents[j] = byte(rand.Intn(94) + 33)
		}
		_, err = f.Write(f_contents)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		// temporary code until we decide what to do with the MD5.
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		MD5 := strconv.FormatInt(r.Int63(), 10)
		file := model.File{
			Name: f.Name(),
			Size: int64(len(f_contents)),
			MD5:  MD5,
		}

		file.BoreholeIDs = append(file.BoreholeIDs, 1)
		file.InventoryIDs = append(file.InventoryIDs, 1)
		file.OutcropIDs = append(file.OutcropIDs, 1)
		file.ProspectIDs = append(file.ProspectIDs, 1)
		file.WellIDs = append(file.WellIDs, 1)

		var fs_file fsutil.File
		err = db.PutFile(&file, func() error {
			mt := mime.TypeByExtension(filepath.Ext(f.Name()))
			if mt == "" {
				mt = "application/octet-stream"
			}
			fs_file = fsutil.File{
				Name:         fmt.Sprintf("%d/%s", file.ID, file_name),
				Size:         f_content_bytes,
				LastModified: time.Now(),
				ContentType:  mt,
				Content:      f,
			}
			//Add the file to the filestore
			err = fs.PutFile(&fs_file)
			if err != nil {
				return fmt.Errorf("error putting file in filestore: %w", err)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("FilePut executed successfully")
		}
		err = db.DeleteFile(&file, true)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("DeleteFile executed successfully:", file.ID)
			err = fs.DeleteFile(&fs_file)
			if err != nil {
				t.Fatal(err)
			} else {
				t.Log("file deleted from filestore:", fs_file.Name)
			}
		}
	}
}

func TestFsPutGetFile(t *testing.T) {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		t.Fatal("Error:", err)
	}
	cfg, err := config.Load(filepath.Join(config_dir, "gmc.test.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	fs, err := filestore.New(cfg.FileStore)
	if err != nil {
		t.Fatal(err)
	}
	params := []params{
		{P_type: "borehole", Content_bytes: 0, ID: 1},
		{P_type: "inventory", Content_bytes: 100},
		{P_type: "outcrop", Content_bytes: 1000},
		{P_type: "prospect", Content_bytes: 10000},
		{P_type: "well", Content_bytes: 100000},
		{P_type: "", Content_bytes: 100, ID: 1},
	}

	for _, p := range params {
		f, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatal(err)
		}
		file_name := filepath.Base(f.Name())
		defer os.Remove(f.Name())
		defer f.Close()
		f_content_bytes := p.Content_bytes
		f_contents := make([]byte, f_content_bytes)
		for j := range f_contents {
			f_contents[j] = byte(rand.Intn(94) + 33)
		}
		_, err = f.Write(f_contents)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		// temporary code until we decide what to do with the MD5.
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		MD5 := strconv.FormatInt(r.Int63(), 10)
		file := model.File{
			Name: f.Name(),
			Size: int64(len(f_contents)),
			MD5:  MD5,
		}
		// if id is 0, create a random by_file_id
		if p.ID == 0 {
			file.ID = r.Int31n(5000) + 1
		} else {
			file.ID = p.ID
		}
		switch p.P_type {
		case "borehole":
			file.BoreholeIDs = append(file.BoreholeIDs, 1)
		case "inventory":
			file.InventoryIDs = append(file.InventoryIDs, 1)
		case "outcrop":
			file.OutcropIDs = append(file.OutcropIDs, 1)
		case "prospect":
			file.ProspectIDs = append(file.ProspectIDs, 1)
		case "well":
			file.WellIDs = append(file.WellIDs, 1)
		}
		var fs_file fsutil.File
		mt := mime.TypeByExtension(filepath.Ext(f.Name()))
		if mt == "" {
			mt = "application/octet-stream"
		}
		fs_file = fsutil.File{
			Name:         fmt.Sprintf("%d/%s", file.ID, file_name),
			Size:         f_content_bytes,
			LastModified: time.Now(),
			ContentType:  mt,
			Content:      f,
		}
		err = fs.PutFile(&fs_file)
		if err != nil {
			t.Fatal("failure: file not added to the filestore", file.ID)
		} else {
			t.Log("success: file added to the filestore", file.ID)
		}

		retrieved_file, err := fs.GetFile(strconv.Itoa(int(file.ID)))
		if err != nil {
			t.Fatal("failure: unable to get file from the filestore", file.ID, err)
		} else {
			t.Log("success: file retrieved from the filestore", retrieved_file.Name)
		}

		err = fs.DeleteFile(&fs_file)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("file deleted from filestore:", fs_file.Name)
		}
	}
}

func TestDBPutGetFile(t *testing.T) {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		t.Fatal("Error:", err)
	}
	cfg, err := config.Load(filepath.Join(config_dir, "gmc.test.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	db, err := db.New(cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	params := []params{
		{P_type: "borehole", Content_bytes: 0},
		{P_type: "inventory", Content_bytes: 100},
		{P_type: "outcrop", Content_bytes: 1000},
		{P_type: "prospect", Content_bytes: 10000},
		{P_type: "well", Content_bytes: 100000},
		{P_type: "", Content_bytes: 100},
	}
	for _, p := range params {
		f, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		defer f.Close()
		f_content_bytes := p.Content_bytes
		f_contents := make([]byte, f_content_bytes)
		for j := range f_contents {
			f_contents[j] = byte(rand.Intn(94) + 33)
		}
		_, err = f.Write(f_contents)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		// temporary code until we decide what to do with the MD5.
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		MD5 := strconv.FormatInt(r.Int63(), 10)
		file := model.File{
			Name: f.Name(),
			Size: int64(len(f_contents)),
			MD5:  MD5,
		}
		switch p.P_type {
		case "borehole":
			file.BoreholeIDs = append(file.BoreholeIDs, 1)
		case "inventory":
			file.InventoryIDs = append(file.InventoryIDs, 1)
		case "outcrop":
			file.OutcropIDs = append(file.OutcropIDs, 1)
		case "prospect":
			file.ProspectIDs = append(file.ProspectIDs, 1)
		case "well":
			file.WellIDs = append(file.WellIDs, 1)
		}
		mt := mime.TypeByExtension(filepath.Ext(f.Name()))
		if mt == "" {
			mt = "application/octet-stream"
		}
		err = db.PutFile(&file, func() error {
			return nil
		})
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("FilePut executed successfully")
		}
		retrieved_file, err := db.GetFile(int(file.ID))
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("FileGet executed successfully", retrieved_file.ID)
		}
		err = db.DeleteFile(&file, true)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("file's metadata deleted from db:", file.Name)
		}
	}
}
