package controller

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"short2/internal/data"
	"short2/internal/server"
	"short2/internal/ui"
)

const cookieName = "jCookShort2"

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
	mux.HandleFunc("GET "+ctx+"/go/{id}", ctrl.getGo)
	mux.HandleFunc("POST "+ctx+"/login", ctrl.postLogin)
	mux.HandleFunc("POST "+ctx+"/", ctrl.postShorten)

	return http.ListenAndServe(cfg.ServerAddress, NotFoundLoggerMiddleware(mux))
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// NotFoundLoggerMiddleware logs requests that result in a 404 status
func NotFoundLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Default to 200 OK if WriteHeader isn't explicitly called
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rec, r)

		// Check if the response resulted in a 404
		if rec.statusCode == http.StatusNotFound {
			log.Printf("[404 NOT FOUND] Method: %s | URL: %s | RemoteAddr: %s", r.Method, r.URL.Path, r.RemoteAddr)
		}
	})
}

func (c *Controller) processToken(r *http.Request) (*ui.ViewHeader, string, error) {
	// get the cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, "", errors.New("Not logged in")
		} else {
			return nil, "", errors.New("Error while retrieving the cookie")
		}
	}

	token := cookie.Value

	servr := server.NewServer(c.configData, c.db)

	userId, err := servr.GetUserIdByTokenAndAdvance(token)

	if err != nil {
		return nil, "", err
	}

	person, err := servr.GetUserById(userId)
	if err != nil {
		return nil, "", err
	}

	return &ui.ViewHeader{
		Login:    person.Login,
		FullName: person.FullName,
		Title:    "",
		Context:  c.configData.Context,
		Error:    "",
		UserID:   person.ID,
	}, token, nil
}
