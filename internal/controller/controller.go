package controller

import (
	"database/sql"
	"html/template"
	"net/http"
	"short2/internal/data"
	"short2/internal/ui"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func StartServer(cfg *data.ConfigData, db *sql.DB, templates map[string]*template.Template) error {
	ctx := cfg.Context
	ctrl := NewController()

	mux := http.NewServeMux()
	mux.Handle("GET "+ctx+"/static/", http.StripPrefix(ctx, ui.StaticHandler()))
	mux.HandleFunc("GET "+ctx+"/", ctrl.getRoot)

	return http.ListenAndServe(cfg.ServerAddress, mux)
}
