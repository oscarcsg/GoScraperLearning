package logging

type ExternalLogConfig struct {
	Provider          string
	MinLevel          string
	TelChatId         string
	TelBotToken       string
	WebhookURL        string
	WebhookAuthHeader string
}

type LogConfig struct {
	FilePath          string
	GlobalMinLevel    string
	FileMaxSize       uint16
	FileMaxAge        uint16
	FileMaxBackups    uint8
	FileCompress      bool
	Terminal          bool
}

type alertPayload struct {
	Level   Level
	Message string
	Fields  []Field
}

type minLogLevels struct {
	MinGlobalLogLevel   Level
	MinExternalLogLevel Level
}



// ------ Log Levels ------ //

type Level uint8

const (
	NONE  Level = 0
	DEBUG Level = 1
	INFO  Level = 2
	WARN  Level = 3
	ERROR Level = 4
	FATAL Level = 5
)