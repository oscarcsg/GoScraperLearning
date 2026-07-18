package main

import (
	"fmt"
	"go-scraper-learning/internal/config"
	"go-scraper-learning/internal/logging"
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

	logging.Init(appConfig.Logs, appConfig.ExtLogs)
	defer logging.Close()

	initialMsg := fmt.Sprintf(
		"------ %s - Launched in Version %s ------",
		Name,
		Version,
	)

	logging.Always(initialMsg)
}