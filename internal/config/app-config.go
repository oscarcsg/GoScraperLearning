package config

type AppConfig struct {
	Logs     LogConfig
	TelLogs  TelegramLogConfig
	DBEngine string
	DBString string
}