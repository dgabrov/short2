package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"short2/internal/controller"
	"short2/internal/data"
	"short2/internal/ui"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func Start() error {
	cfg, err := loadConfigData()
	if err != nil {
		return err
	}

	db, err := connectToDb(cfg.Db)
	if err != nil {
		return err
	}

	// parse the templates to have them
	templates, err := ui.ParseItems()
	if err != nil {
		return err
	}

	// let's start the http service
	err = controller.StartServer(cfg, db, templates)

	defer db.Close()

	return err
}

func connectToDb(dbconfig *data.DbConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbconfig.Login, dbconfig.Password, dbconfig.Machine, dbconfig.Port, dbconfig.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadConfigData() (*data.ConfigData, error) {
	configFileName := os.Getenv("CONFIG_FILE")
	if len(strings.TrimSpace(configFileName)) == 0 {
		return nil, errors.New("CONFIG_FILE variable not found")
	}

	configFileContents, err := os.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	var cfg data.ConfigData
	err = json.Unmarshal(configFileContents, &cfg)
	if err != nil {
		return nil, err
	}

	// ok log the contents
	slog.Info(fmt.Sprintf("config - serverAddress - %s", cfg.ServerAddress))
	slog.Info(fmt.Sprintf("config - context - %s", cfg.Context))
	slog.Info(fmt.Sprintf("config - db - database - %s", cfg.Db.Database))
	slog.Info(fmt.Sprintf("config - db - port - %d", cfg.Db.Port))
	slog.Info(fmt.Sprintf("config - db - login - %s", cfg.Db.Login))
	slog.Info(fmt.Sprintf("config - db - machine - %s", cfg.Db.Machine))

	return &cfg, nil
}
