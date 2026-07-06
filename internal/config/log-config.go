package config

type LogConfig struct {
	FilePath       string
	FileMaxSize    uint16
	FileMaxAge     uint16
	FileMaxBackups uint8
	FileCompress   bool
	Terminal       bool
}