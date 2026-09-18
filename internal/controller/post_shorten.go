package controller

import (
	"fmt"
	"net/http"
	"short2/internal/server"
	"short2/internal/ui"
	"strings"
)

func (c *Controller) postShorten(w http.ResponseWriter, r *http.Request) {
	header, _, err := c.processToken(r)
	if err != nil {
		_ = ui.RenderLogin(w, c.templates, &ui.LoginView{
			ViewHeader: ui.ViewHeader{},
			Login:      "",
		})

		return
	}

	userID := header.UserID

	err = r.ParseForm()
	if err != nil {
		_ = ui.RenderAddUrl(w, c.templates, &ui.AddUrlView{
			ViewHeader: *header,
			LongUrl:    "",
			ShowShort:  false,
			ShortUrl:   "",
		})

		return
	}

	longUrl := r.FormValue("longUrl")
	longUrl = strings.TrimSpace(longUrl)

	servr := server.NewServer(c.configData, c.db)
	shortUrl, err := servr.ShortenAndSave(userID, longUrl)

	shortUrl = fmt.Sprintf("%s%s/%s", c.configData.ShortUrlPrefix, c.configData.Context, shortUrl)

	if err != nil {
		_ = ui.RenderAddUrl(w, c.templates, &ui.AddUrlView{
			ViewHeader: *header,
			LongUrl:    longUrl,
			ShowShort:  true,
			ShortUrl:   shortUrl,
		})
	}
}
