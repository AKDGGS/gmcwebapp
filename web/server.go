package web

import (
	"context"
	"gmc/assets"
	"gmc/auth"
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
	Auths     auth.Auths
	AssetPath string
	http      http.Server
}

func (srv *Server) Start() error {
	assets.Initialize(srv.AssetPath)

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
