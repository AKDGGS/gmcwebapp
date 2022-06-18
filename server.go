package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gmc/auth"
	"gmc/config"
	"gmc/db"
	"gmc/filestore"
	"gmc/web"
)

func serverCommand(cfg *config.Config, exec string) {
	db, err := db.New(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		os.Exit(1)
	}

	auths, err := auth.NewAuths(cfg.SessionKeyBytes(), cfg.MaxAge, cfg.Auths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		os.Exit(1)
	}

	stor, err := filestore.New(cfg.FileStore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		os.Exit(1)
	}

	srv := web.Server{Config: cfg, DB: db, FileStore: stor, Auths: auths}
	if cfg.AutoShutdown {
		expath, err := filepath.Abs(exec)
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
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		os.Exit(1)
	}
}
