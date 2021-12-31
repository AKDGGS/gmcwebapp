package web

import (
	"context"
	"gmc/assets"
	"gmc/config"
	"gmc/db"
	"gmc/filestore"
	"net"
	"net/http"
)

type Server struct {
	Config    *config.Config
	DB        db.DB
	FileStore filestore.FileStore
	http      http.Server
}

func (srv *Server) Start() error {
	// Build asset cache
	if err := assets.Initialize(); err != nil {
		return err
	}

	listen, err := net.Listen("tcp", srv.Config.ListenAddress)
	if err != nil {
		return err
	}

	srv.http = http.Server{Handler: srv}
	err = srv.http.Serve(listen)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (srv *Server) Shutdown() {
	srv.http.Shutdown(context.Background())
}
