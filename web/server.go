package web

import (
	"context"
	"gmc/assets"
	"gmc/auth"
	"gmc/auth/user"
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
	Auths     []auth.Auth
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

func (srv *Server) AuthRequired(w http.ResponseWriter, r *http.Request) (*user.User, error) {
	for _, au := range srv.Auths {
		user, err := au.AuthRequired(w, r)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}
	return nil, nil
}

func (srv *Server) AuthOptional(w http.ResponseWriter, r *http.Request) (*user.User, error) {
	for _, au := range srv.Auths {
		user, err := au.AuthOptional(w, r)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}
	return nil, nil
}
