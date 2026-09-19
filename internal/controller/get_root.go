package controller

import (
	"net/http"
	"short2/internal/ui"
)

func (c *Controller) getRoot(w http.ResponseWriter, r *http.Request) {
	header, _, err := c.processToken(r)
	if err != nil {
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

		_ = ui.RenderLogin(w, c.templates, &data)
	} else {
		header.Title = "Shortener - Add Url"
		data := ui.AddUrlView{
			ViewHeader: *header,
			LongUrl:    "",
			ShowShort:  false,
			ShortUrl:   "",
		}

		_ = ui.RenderAddUrl(w, c.templates, &data)
	}

}
