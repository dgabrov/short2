package ui

import (
	"embed"
	"net/http"
)

//go:embed templates
var templateFiles embed.FS

//go:embed static
var staticFiles embed.FS

// StaticHandler serves embedded static assets (CSS, etc.) under /static/.
func StaticHandler() http.Handler {
	fileSystem := http.FS(staticFiles)

	return http.FileServer(fileSystem)
}
