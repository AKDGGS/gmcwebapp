package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gmc/auth"
	"gmc/config"
	"gmc/db"
	"gmc/filestore"
	"gmc/search"
	"gmc/web"
)

func ServerCommand(cfg *config.Config, exec string) int {
	db, err := db.New(cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		return 1
	}

	sea, err := search.New(cfg.Search)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		return 1
	}

	auths, err := auth.NewAuths(cfg.SessionKeyBytes(), cfg.MaxAge, cfg.Auths, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		return 1
	}

	stor, err := filestore.New(cfg.FileStore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		return 1
	}

	srv := web.Server{
		Config: cfg, DB: db, Search: sea, FileStore: stor, Auths: auths,
	}
	if cfg.AutoShutdown {
		expath, err := filepath.Abs(exec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "absolute path error: %s\n", err.Error())
			return 1
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
		return 1
	}

	return 0
}
