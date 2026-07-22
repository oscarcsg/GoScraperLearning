package main

import (
	"fmt"
	"go-scraper-learning/internal/config"
	"go-scraper-learning/internal/database"
	"go-scraper-learning/internal/logging"
	"go-scraper-learning/internal/scrapper"
)

var (
	Name        = "unknown"
	Author      = "unknown"
	Description = "unknown"
	Version     = "unknown"
	BuildTime   = "unknown"
	CommitID    = "unknown"
)

func main() {
	appConfig := config.LoadConfig()

	// --- LOGGING INITIALIZATION --- //
	logging.Init(appConfig.Logs, appConfig.ExtLogs)
	defer logging.Close()

	// --- LAUNCH MESSAGE --- //
	initialMsg := fmt.Sprintf(
		"------ %s - Launched in Version %s ------",
		Name,
		Version,
	)

	logging.Always(initialMsg)

	// --- DATABASE INITIALIZATION --- //
	db, err := database.Init(appConfig.DBConfig)
	if err != nil {
		logging.Fatal(
			"Database initialization failed.",
			logging.ErrorType(err),
		)
	}
	defer db.Close()
	if appConfig.DBConfig.DBEngine == "sqlite" {
		scrapper.CreateDatabaseSQLite(db)
	}

	// --- SCRAPPER --- //
	scrapper.Init(
		"https://books.toscrape.com/",
		"historial/historial-dev.txt",
	)
}