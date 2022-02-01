package main

import (
	"flag"
	"fmt"
	"gmc/config"
	"gmc/db"
	"gmc/filestore"
	"gmc/web"
	"os"
	"path/filepath"
	"time"
)

func serverCommand() {
	cmd := flag.NewFlagSet("server", flag.ExitOnError)
	cmd.SetOutput(os.Stdout)
	cpath := cmd.String("conf", "", "path to configuration")
	autos := cmd.Bool("s", false, "automatic shutdown on executable change")
	assets := cmd.String("assets", "", "override embedded assets with assets from path")
	cmd.Parse(os.Args[2:])

	cfg, err := config.Load(*cpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}

	db, err := db.New(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}

	stor, err := filestore.New(cfg.FileStore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}

	srv := web.Server{Config: cfg, DB: db, FileStore: stor}
	if *assets != "" {
		srv.AssetPath = *assets
	}
	if *autos {
		expath, err := filepath.Abs(os.Args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "absolute path error: %s\n", err.Error())
			os.Exit(1)
		}

		go func() {
			var t time.Time

			for {
				time.Sleep(time.Second)

				fi, err := os.Stat(expath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "stat error: %s\n", err.Error())
					continue
				}

				if t.IsZero() {
					t = fi.ModTime()
				} else if fi.ModTime().After(t) {
					srv.Shutdown()
				}
			}
		}()
	}

	if err = srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}
}
