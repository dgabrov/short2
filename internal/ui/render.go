package ui

import (
	"embed"
	"html/template"
	"io"
	"net/http"
	"strings"
)

//go:embed templates
var templateFiles embed.FS

//go:embed static
var staticFiles embed.FS

const LoginTemplate = "login"
const AddUrlTemplate = "addUrl"

// StaticHandler serves embedded static assets (CSS, etc.) under /static/.
func StaticHandler() http.Handler {
	fileSystem := http.FS(staticFiles)

	return http.FileServer(fileSystem)
}

type ViewHeader struct {
	Login    string
	FullName string
	Title    string
	Context  string
	Error    string
	UserID   string
}

func (c *ViewHeader) IsLoggedIn() bool {
	trimmed := strings.TrimSpace(c.Login)

	return len(trimmed) > 0
}

type LoginView struct {
	ViewHeader
	Login string
}

type AddUrlView struct {
	ViewHeader
	LongUrl   string
	ShowShort bool
	ShortUrl  string
}

func ParseItems() (map[string]*template.Template, error) {
	templateMap := make(map[string]*template.Template)

	// login
	loginTemplate, err := template.ParseFS(templateFiles, "templates/header.html", "templates/login.html")
	if err != nil {
		return nil, err
	}

	// addUrl
	addUrlTemplate, err := template.ParseFS(templateFiles, "templates/header.html", "templates/addurl.html")
	if err != nil {
		return nil, err
	}

	templateMap[LoginTemplate] = loginTemplate
	templateMap[AddUrlTemplate] = addUrlTemplate

	return templateMap, nil
}

func RenderLogin(writer io.Writer, templates map[string]*template.Template, data *LoginView) error {
	templ := templates[LoginTemplate]

	return templ.ExecuteTemplate(writer, "header", data)
}

func RenderAddUrl(writer io.Writer, templates map[string]*template.Template, data *AddUrlView) error {
	templ := templates[AddUrlTemplate]

	return templ.ExecuteTemplate(writer, "header", data)
}
