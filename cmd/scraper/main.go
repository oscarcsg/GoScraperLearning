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
	fmt.Println("Learn scrapping project initialization.") // Log

	appConfig := config.LoadConfig()
}