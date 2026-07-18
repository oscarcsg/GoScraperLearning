package config

import "go-scraper-learning/internal/logging"

type AppConfig struct {
	ExtLogs  logging.ExternalLogConfig
	Logs     logging.LogConfig
	DBEngine string
	DBString string
}