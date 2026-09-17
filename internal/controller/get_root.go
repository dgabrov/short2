package controller

import (
	"log/slog"
	"net/http"
	"short2/internal/ui"
)

func (c *Controller) getRoot(w http.ResponseWriter, r *http.Request) {
	data := ui.AddUrlView{
		ViewHeader: ui.ViewHeader{
			Login:    "",
			FullName: "Buben Lampa",
			Title:    "De title here",
			Context:  "short",
			Error:    "the error message should stay here corrected",
		},
		LongUrl:   "laila",
		ShowShort: true,
		ShortUrl:  "http://lampa.org?here",
	}

	err := ui.RenderAddUrl(w, c.templates, &data)

	if err != nil {
		slog.Info(err.Error())
	}

}
