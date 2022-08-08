package web

import (
	"context"
	"net"
	"net/http"

	"gmc/auth"
	"gmc/config"
	"gmc/db"
	"gmc/filestore"
)

type Server struct {
	Config    *config.Config
	DB        db.DB
	FileStore filestore.FileStore
	Auths     *auth.Auths
	http      http.Server
}

func (srv *Server) Start() error {
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
