package controller

import (
	"database/sql"
	"html/template"
	"net/http"
	"short2/internal/data"
	"short2/internal/ui"
)

type Controller struct {
	configData *data.ConfigData
	db         *sql.DB
	templates  map[string]*template.Template
}

func NewController(configData *data.ConfigData, db *sql.DB, templates map[string]*template.Template) *Controller {
	return &Controller{
		configData: configData,
		db:         db,
		templates:  templates,
	}
}

func StartServer(cfg *data.ConfigData, db *sql.DB, templates map[string]*template.Template) error {
	ctx := cfg.Context
	ctrl := NewController(cfg, db, templates)

	mux := http.NewServeMux()
	mux.Handle("GET "+ctx+"/static/", http.StripPrefix(ctx, ui.StaticHandler()))
	mux.HandleFunc("GET "+ctx+"/", ctrl.getRoot)
	mux.HandleFunc("GET "+ctx+"/logout", ctrl.getLogout)
	mux.HandleFunc("GET "+ctx+"/go", ctrl.getGo)
	mux.HandleFunc("POST "+ctx+"/login", ctrl.postLogin)
	mux.HandleFunc("POST "+ctx+"/", ctrl.postShorten)

	return http.ListenAndServe(cfg.ServerAddress, mux)
}
