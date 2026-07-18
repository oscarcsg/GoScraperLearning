// internal/logging/logger.go
package logging

import (
	"fmt"
	"go-scraper-learning/internal/util"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLog *zap.Logger
var alwaysLog *zap.Logger

var minLevels minLogLevels = minLogLevels{}

func Init(logConfig LocalLogConfig, extLogConfig ExternalLogConfig) {
	setMinimumLogLevels(&minLevels, &logConfig, &extLogConfig)

	// If the min levels are NONE globally or each one is NONE, just do not activate the system
	if minLevels.TerminalLogMinLevel == NONE && 
	   minLevels.FileLogMinLevel     == NONE &&
	   minLevels.ExternalLogMinLevel == NONE {

		zapLog = zap.NewNop()
		return
	}

	// Cores (where it will log)
	var cores []zapcore.Core
	var alwaysCores []zapcore.Core

	// Configure terminal
	if logConfig.Terminal && minLevels.TerminalLogMinLevel != NONE {
		consoleConfig := zap.NewDevelopmentEncoderConfig()
		consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
		
		consoleWriter := zapcore.AddSync(os.Stdout)
		cores = append(
			cores,
			zapcore.NewCore(
				consoleEncoder,
				consoleWriter,
				toZapLevel(minLevels.TerminalLogMinLevel),
			),
		)
		// Ignore the Level and logs EVERYTHING
		alwaysCores = append(alwaysCores, zapcore.NewCore(consoleEncoder, consoleWriter, zap.DebugLevel))
	}

	// Configure file
	if minLevels.FileLogMinLevel != NONE && logConfig.FilePath != "" {
		// Configure lumberjack
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logConfig.FilePath,
			MaxSize:    int(logConfig.FileMaxSize),
			MaxAge:     int(logConfig.FileMaxAge),
			MaxBackups: int(logConfig.FileMaxBackups),
			Compress:   logConfig.FileCompress,
		})

		var fileEncoder zapcore.Encoder

		// Get the file extension to decide which format to use
		fileExtension := strings.Split(strings.ToLower(logConfig.FilePath), ".")
		if fileExtension[len(fileExtension)-1] == "json" {
			// Configure JSON encoder
			fileEncoderConfig := zap.NewProductionEncoderConfig()
			fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			fileEncoder = zapcore.NewJSONEncoder(fileEncoderConfig)
		} else {
			fileEncoderConfig := zap.NewDevelopmentEncoderConfig()
			fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			fileEncoder = zapcore.NewConsoleEncoder(fileEncoderConfig)
		}

		cores = append(
			cores,
			zapcore.NewCore(
				fileEncoder,
				fileWriter,
				toZapLevel(minLevels.FileLogMinLevel),
			),
		)

		// Ignore the Level and logs EVERYTHING
		alwaysCores = append(
			alwaysCores,
			zapcore.NewCore(
				fileEncoder,
				fileWriter,
				zap.DebugLevel,
			),
		)
	}

	// Initialize async worker for external logging
	if minLevels.ExternalLogMinLevel != NONE {
		startAlertWorker(extLogConfig)
	}

	// Create the instances
	// AddCallerSkip(1) will write the real file that made the log, not this file
	if len(cores) > 0 {
		zapLog = zap.New(zapcore.NewTee(cores...), zap.AddCallerSkip(1))
	} else {
		zapLog = zap.NewNop()
	}

	if len(alwaysCores) > 0 {
		alwaysLog = zap.New(zapcore.NewTee(alwaysCores...), zap.AddCallerSkip(1))
	} else {
		alwaysLog = zap.NewNop()
	}
}

// Call Close using defer in main()

func Close() {
	if zapLog != nil {
		_ = zapLog.Sync()
	}
	if alwaysLog != nil {
		_ = alwaysLog.Sync()
	}

	// Give 5 seconds so the logging system can empty as many logs as it can from the alertChannel
	time.Sleep(5 * time.Second)
}



// --------------------------------- //
//            OPERATIONS             //
// --------------------------------- //

func Always(msg string, fields ...Field) {
	alwaysLog.Info(msg, fields...)
	queueAlert(ALWAYS, msg, fields)
}

func Debug(msg string, fields ...Field) {
	zapLog.Debug(msg, fields...)
	if DEBUG >= minLevels.ExternalLogMinLevel {
		queueAlert(DEBUG, msg, fields)
	}
}

func Info(msg string, fields ...Field) {
	zapLog.Info(msg, fields...)
	if INFO >= minLevels.ExternalLogMinLevel {
		queueAlert(INFO, msg, fields)
	}
}

func Warn(msg string, fields ...Field) {
	zapLog.Warn(msg, fields...)
	if WARN >= minLevels.ExternalLogMinLevel {
		queueAlert(WARN, msg, fields)
	}
}

func Error(msg string, fields ...Field) {
	zapLog.Error(msg, fields...)
	if ERROR >= minLevels.ExternalLogMinLevel {
		queueAlert(ERROR, msg, fields)
	}
}

func Fatal(msg string, fields ...Field) {
	if FATAL >= minLevels.ExternalLogMinLevel {
		queueAlert(FATAL, msg, fields)
		// Give margin to send the log before the program exits
		_ = zapLog.Sync()
	}
	zapLog.Fatal(msg, fields...)
}



// Set the minimum levels

func setMinimumLogLevels(minLevels *minLogLevels, logConfig *LocalLogConfig, extLogConfig *ExternalLogConfig) {
	globalMinLevel   := util.TrimToLowerString(logConfig.GlobalMinLevel)
	terminalMinLevel := util.TrimToLowerString(logConfig.TerminalMinLevel)
	fileMinLevel     := util.TrimToLowerString(logConfig.FileMinLevel)
	externalMinLevel := util.TrimToLowerString(extLogConfig.MinLevel)

	// Global
	switch globalMinLevel {
	case "", "d", "default":
		minLevels.GlobalLogMinLevel = WARN
	case "none":
		minLevels.GlobalLogMinLevel = NONE
	case "debug":
		minLevels.GlobalLogMinLevel = DEBUG
	case "info":
		minLevels.GlobalLogMinLevel = INFO
	case "warn":
		minLevels.GlobalLogMinLevel = WARN
	case "error":
		minLevels.GlobalLogMinLevel = ERROR
	case "fatal":
		minLevels.GlobalLogMinLevel = FATAL
	default:
		minLevels.GlobalLogMinLevel = WARN
		fmt.Printf("Value of minimum global log level received is not valid, min global log level is set automatically to default (WARN). [%s]\n", strings.TrimSpace(logConfig.GlobalMinLevel))
	}

	// Terminal
	switch terminalMinLevel {
	case "", "d", "default":
		minLevels.TerminalLogMinLevel = minLevels.GlobalLogMinLevel
	case "none":
		minLevels.TerminalLogMinLevel = NONE
	case "debug":
		minLevels.TerminalLogMinLevel = DEBUG
	case "info":
		minLevels.TerminalLogMinLevel = INFO
	case "warn":
		minLevels.TerminalLogMinLevel = WARN
	case "error":
		minLevels.TerminalLogMinLevel = ERROR
	case "fatal":
		minLevels.TerminalLogMinLevel = FATAL
	default:
		minLevels.TerminalLogMinLevel = WARN
		fmt.Printf("Value of minimum terminal log level received is not valid, min terminal log level is set automatically to default (WARN). [%s]\n", strings.TrimSpace(logConfig.TerminalMinLevel))
	}

	// File
	switch fileMinLevel {
	case "", "d", "default":
		minLevels.FileLogMinLevel = minLevels.GlobalLogMinLevel
	case "none":
		minLevels.FileLogMinLevel = NONE
	case "debug":
		minLevels.FileLogMinLevel = DEBUG
	case "info":
		minLevels.FileLogMinLevel = INFO
	case "warn":
		minLevels.FileLogMinLevel = WARN
	case "error":
		minLevels.FileLogMinLevel = ERROR
	case "fatal":
		minLevels.FileLogMinLevel = FATAL
	default:
		minLevels.FileLogMinLevel = WARN
		fmt.Printf("Value of minimum file log level received is not valid, min file log level is set automatically to default (WARN). [%s]\n", strings.TrimSpace(logConfig.FileMinLevel))
	}

	// External
	switch externalMinLevel {
	case "", "d", "default":
		minLevels.ExternalLogMinLevel = minLevels.GlobalLogMinLevel
	case "none":
		minLevels.ExternalLogMinLevel = NONE
	case "debug":
		minLevels.ExternalLogMinLevel = DEBUG
	case "info":
		minLevels.ExternalLogMinLevel = INFO
	case "warn":
		minLevels.ExternalLogMinLevel = WARN
	case "error":
		minLevels.ExternalLogMinLevel = ERROR
	case "fatal":
		minLevels.ExternalLogMinLevel = FATAL
	default:
		minLevels.ExternalLogMinLevel = WARN
		fmt.Printf("Value of minimum external log level received is not valid, min external log level is set automatically to default (WARN). [%s]\n", strings.TrimSpace(extLogConfig.MinLevel))
	}
}

func toZapLevel(lvl Level) (zapcore.LevelEnabler) {
	switch lvl {
	case DEBUG: return zapcore.DebugLevel
	case INFO:  return zapcore.InfoLevel
	case WARN:  return zapcore.WarnLevel
	case ERROR: return zapcore.ErrorLevel
	case FATAL: return zapcore.FatalLevel
	default:    return zapcore.WarnLevel
	}
}
