package controller

import (
	"net/http"
	"short2/internal/server"
	"strings"
)

func (c *Controller) getGo(w http.ResponseWriter, r *http.Request) {
	// look for the  entry
	shortID := r.PathValue("id")
	shortID = strings.TrimSpace(shortID)
	if shortID == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
	}

	// get the entry
	servr := server.NewServer(c.configData, c.db)
	url, err := servr.GetUrlByShortID(shortID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	http.Redirect(w, r, url, http.StatusFound)
}
