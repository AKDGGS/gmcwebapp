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
	webu "gmc/web/util"
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

	// Main page redirect
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		webu.Redirect(w, "inventory/search", http.StatusFound)
	})

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
	mux.HandleFunc("/api/audit.json", srv.ServeAPIAudit)
	mux.HandleFunc("/api/move_contents.json", srv.ServeAPIMoveInventoryAndContainersContents)
	mux.HandleFunc("/api/add_container.json", srv.ServeAPIAddContainer)
	mux.HandleFunc("/api/add_inventory.json", srv.ServeAPIAddInventory)
	mux.HandleFunc("/api/add_inventory_quality.json", srv.ServeAPIAddInventoryQuality)
	mux.HandleFunc("/api/recode.json", srv.ServeAPIRecodeInventoryAndContainer)

	// Serves list of quality issues
	mux.HandleFunc("/issues.json", srv.ServeIssues)
	// Serves list of keywords
	mux.HandleFunc("/keywords.json", srv.ServeKeywords)
	// Serves list of collections
	mux.HandleFunc("/collections.json", srv.ServeCollections)
	// Serves list of prospects
	mux.HandleFunc("/prospects.json", srv.ServeProspects)

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
	mux.HandleFunc("/inventory/stash.json", srv.ServeInventoryStash)
	mux.HandleFunc("/inventory/search.json", srv.ServeSearchInventoryJSON)
	mux.HandleFunc("/inventory/search.csv", srv.ServeSearchInventoryCSV)
	mux.HandleFunc("/inventory/search.pdf", srv.ServeSearchInventoryPDF)
	mux.HandleFunc("/inventory/search.geojson", srv.ServeSearchInventoryGeoJSON)
	mux.HandleFunc("/inventory/search", srv.ServeSearchInventoryPage)
	mux.HandleFunc("/inventory/search-help", srv.ServeSearchInventoryHelp)

	srv.http = http.Server{Handler: mux}

	if srv.Config.ListenCertificate != "" && srv.Config.ListenKey != "" {
		err = srv.http.ServeTLS(listen, srv.Config.ListenCertificate,
			srv.Config.ListenKey)
	} else {
		err = srv.http.Serve(listen)
	}
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (srv *Server) Shutdown() {
	srv.http.Shutdown(context.Background())
}
