package main

import (
	"fmt"
	"go-scraper-learning/internal/config"
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
	fmt.Println("Learn scrapping project initialization.")
	config.LoadConfig()
}