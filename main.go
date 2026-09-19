package main

import (
	"log/slog"
	"short2/internal/app"
)

func main() {
	err := app.Start()

	if err != nil {
		slog.Error(err.Error())
	}
}
