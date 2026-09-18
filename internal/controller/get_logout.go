package controller

import (
	"net/http"
	"short2/internal/server"
	"short2/internal/ui"
	"time"
)

func (c *Controller) getLogout(w http.ResponseWriter, r *http.Request) {
	_, token, _ := c.processToken(r)

	servr := server.NewServer(c.configData, c.db)
	if len(token) > 0 {
		_ = servr.ExpireSessionByToken(token)
	}

	// then render the login
	data := ui.LoginView{
		ViewHeader: ui.ViewHeader{
			Login:    "",
			FullName: "",
			Title:    "Login",
			Context:  c.configData.Context,
			Error:    "",
		},
		Login: "",
	}

	// delete the cookie if possible

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/", // MUST match the Path of the cookie you want to delete
		HttpOnly: true,
		MaxAge:   -1,              // Tells browser to delete immediately
		Expires:  time.Unix(0, 0), // Fallback for older browsers
	})

	_ = ui.RenderLogin(w, c.templates, &data)
}
