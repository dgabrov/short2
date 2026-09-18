package controller

import (
	"net/http"
	"short2/internal/ui"
)

func (c *Controller) getRoot(w http.ResponseWriter, r *http.Request) {
	data := ui.LoginView{
		ViewHeader: ui.ViewHeader{
			Login:    "",
			FullName: "",
			Title:    "",
			Context:  "",
			Error:    "",
		},
		Login: "",
	}

	_ = ui.RenderLogin(w, c.templates, &data)
}
