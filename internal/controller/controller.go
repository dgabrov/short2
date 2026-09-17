package controller

import (
	"database/sql"
	"net/http"
	"short2/internal/data"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func StartServer(cfg *data.ConfigData, db *sql.DB) error {
	ctx := cfg.Context
	ctrl := NewController()

	mux := http.NewServeMux()
	mux.HandleFunc("GET "+ctx+"/", ctrl.getRoot)

	return http.ListenAndServe(cfg.ServerAddress, mux)
}
