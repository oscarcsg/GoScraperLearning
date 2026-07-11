package config

type AppConfig struct {
	ExtLogs  ExternalLogConfig
	Logs     LogConfig
	DBEngine string
	DBString string
}