package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"short2/internal/data"
	"short2/internal/server"
	"short2/internal/ui"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (c *Controller) postLogin(w http.ResponseWriter, r *http.Request) {
	err := c.processPostLogin(w, r)

	if err != nil {
		_ = ui.RenderLogin(w, c.templates, &ui.LoginView{
			ViewHeader: ui.ViewHeader{
				Login:    "",
				FullName: "",
				Title:    "",
				Context:  "",
				Error:    err.Error(),
			},
			Login: "",
		})
	}
}

func (c *Controller) processPostLogin(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	login = strings.TrimSpace(login)
	password = strings.TrimSpace(password)

	if len(login) == 0 {
		return errors.New("Please fill out the login")
	}

	loginData := data.LoginData{login, password}
	loginBytes, err := json.Marshal(loginData)
	if err != nil {
		return err
	}

	res, err := http.Post(c.configData.AuthServerUrl, "application/json", bytes.NewReader(loginBytes))
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		bts, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return errors.New(string(bts))
	}

	// now process the response
	var loginAuthData data.LoginAuthData
	err = json.NewDecoder(res.Body).Decode(&loginAuthData)
	if err != nil {
		return err
	}

	// create another session entry with the new token and then load the info and proceed to the next one
	servr := server.NewServer(c.configData, c.db)
	exists, err := servr.ExistsProvidedUserId(loginAuthData.Id)
	if err != nil {
		return err
	}

	var userId string
	if !exists {
		uu, err := uuid.NewV7()
		userId = uu.String()
		err = servr.CreateUser(userId, loginAuthData.Login, loginAuthData.Name, loginAuthData.Id)
		if err != nil {
			return err
		}
	} else {
		userId, err = servr.GetUserIdByProvidedId(loginAuthData.Id)
		if err != nil {
			return err
		}
	}

	token := generateToken()
	expiryDate := time.Now().Add(time.Duration(c.configData.TokenTtlSecond) * time.Second)
	err = servr.CreateSessionData(userId, expiryDate, token)

	if err != nil {
		return err
	}

	// save the http cookie
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	_ = ui.RenderAddUrl(w, c.templates, &ui.AddUrlView{
		ViewHeader: ui.ViewHeader{
			Login:    loginAuthData.Login,
			FullName: loginAuthData.Name,
			Title:    "Add URL",
			Context:  c.configData.Context,
			Error:    "",
		},
		LongUrl:   "",
		ShowShort: false,
		ShortUrl:  "",
	})

	return nil
}

func generateToken() string {
	return uuid.Must(uuid.NewV7()).String()
}
