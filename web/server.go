package web

import (
	"context"
	"net"
	"net/http"

	"gmc/assets"
	"gmc/auth"
	"gmc/config"
	"gmc/db"
	"gmc/filestore"
	"gmc/search"
)

type Server struct {
	Config    *config.Config
	DB        db.DB
	Search    search.Search
	FileStore filestore.FileStore
	Auths     *auth.Auths
	http      http.Server
}

func (srv *Server) Start() error {
	listen, err := net.Listen("tcp", srv.Config.ListenAddress)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	// Interactive login/logout machinery
	mux.HandleFunc("/login", srv.Auths.CheckForm)
	mux.HandleFunc("/logout", srv.Auths.Logout)

	// QA/QC system
	mux.HandleFunc("/qa/{$}", srv.ServeQA)
	mux.HandleFunc("/qa/qa_count.json", srv.ServeQACount)
	mux.HandleFunc("/qa/qa_run.json", srv.ServeQARun)

	// Wells page
	mux.HandleFunc("/wells/{$}", srv.ServeWells)
	mux.HandleFunc("/wells/points.json", srv.ServeWellsPointsJSON)
	mux.HandleFunc("/wells/detail.json", srv.ServeWellsDetail)

	// API endpoints
	mux.HandleFunc("/api/summary.json", srv.ServeAPISummary)
	mux.HandleFunc("/api/inventory.json", srv.ServeAPIInventoryDetail)
	mux.HandleFunc("/api/move.json", srv.ServeAPIMoveInventoryAndContainers)
	mux.HandleFunc("/api/add_audit.json", srv.ServeAPIAudit)
	mux.HandleFunc("/api/move_contents.json", srv.ServeAPIMoveInventoryAndContainersContents)
	mux.HandleFunc("/api/add_container.json", srv.ServeAPIAddContainer)
	mux.HandleFunc("/api/add_inventory.json", srv.ServeAPIAddInventory)
	mux.HandleFunc("/api/add_inventory_quality.json", srv.ServeAPIAddInventoryQuality)
	mux.HandleFunc("/api/recode.json", srv.ServeAPIRecodeInventoryAndContainer)

	// Serves the list of quality issues
	mux.HandleFunc("/issues.json", srv.ServeIssues)

	// Upload page/endpoint
	mux.HandleFunc("/upload/{$}", srv.ServeUpload)

	// General assets
	mux.HandleFunc("/css/", assets.ServeStatic)
	mux.HandleFunc("/js/", assets.ServeStatic)
	mux.HandleFunc("/ol/", assets.ServeStatic)
	mux.HandleFunc("/img/", assets.ServeStatic)

	// Secondary endpoints to assets
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		assets.ServeStaticPath("img/favicon.ico", w, r.Header)
	})
	mux.HandleFunc("/ol/ol.css", func(w http.ResponseWriter, r *http.Request) {
		assets.ServeStaticPath("ol/ol-v7.1.0.css", w, r.Header)
	})
	mux.HandleFunc("/ol/ol.js", func(w http.ResponseWriter, r *http.Request) {
		assets.ServeStaticPath("ol/ol-v7.1.0.js", w, r.Header)
	})
	mux.HandleFunc("/js/mustache.js", func(w http.ResponseWriter, r *http.Request) {
		assets.ServeStaticPath("js/mustache-v4.2.0.js", w, r.Header)
	})

	// Pages for each type by ID
	mux.HandleFunc("/well/{id}", srv.ServeWell)
	mux.HandleFunc("/file/{id}", srv.ServeFile)
	mux.HandleFunc("/prospect/{id}", srv.ServeProspect)
	mux.HandleFunc("/borehole/{id}", srv.ServeBorehole)
	mux.HandleFunc("/outcrop/{id}", srv.ServeOutcrop)
	mux.HandleFunc("/shotline/{id}", srv.ServeShotline)
	mux.HandleFunc("/inventory/{id}", srv.ServeInventory)

	// Used by inventory
	mux.HandleFunc("/stash.json", srv.ServeStash)

	srv.http = http.Server{Handler: mux}
	err = srv.http.Serve(listen)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (srv *Server) Shutdown() {
	srv.http.Shutdown(context.Background())
}
