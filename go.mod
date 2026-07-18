module go-scraper-learning

go 1.26.1

// Environment variables
require github.com/joho/godotenv v1.5.1

// Logging
require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.28.0
	
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)
