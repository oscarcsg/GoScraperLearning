// internal/logging/logger.go
package logging

import (
	"fmt"
	"go-scraper-learning/internal/util"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLog *zap.Logger
var minLevels minLogLevels = minLogLevels{MinGlobalLogLevel: NONE, MinExternalLogLevel: NONE}

func Init(logConfig LogConfig, extLogConfig ExternalLogConfig) {
	setMinimumLogLevels(&minLevels, &logConfig, &extLogConfig)

	// 1. Iniciar worker asíncrono
	startAlertWorker(extLogConfig)

	// 2. Configurar rotación de archivos (Lumberjack)
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.FilePath,
		MaxSize:    int(logConfig.FileMaxSize),
		MaxAge:     int(logConfig.FileMaxAge),
		MaxBackups: int(logConfig.FileMaxBackups),
		Compress:   logConfig.FileCompress,
	})

	// 3. Encoder JSON para el archivo (Máximo rendimiento de lectura automatizada)
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	// 4. Multiplexor de Cores
	var cores []zapcore.Core
	cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, zap.DebugLevel))

	if logConfig.Terminal {
		consoleConfig := zap.NewDevelopmentEncoderConfig()
		consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
		
		consoleWriter := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriter, zap.DebugLevel))
	}

	// 5. Instanciar Zap. AddCallerSkip(1) es vital para que Zap imprima 
	// el archivo real donde ocurrió el error, no este logger.go.
	zapLog = zap.New(zapcore.NewTee(cores...), zap.AddCallerSkip(1))
}



// --------------------------------- //
//            OPERATIONS             //
// --------------------------------- //

func Debug(msg string, fields ...Field) {
	if minLevels.MinGlobalLogLevel >= DEBUG {
		zapLog.Debug(msg, fields...)
	}
	if minLevels.MinGlobalLogLevel >= DEBUG {
		queueAlert(DEBUG, msg, fields)
	}
}

func Info(msg string, fields ...Field) {
	if minLevels.MinGlobalLogLevel >= INFO {
		zapLog.Info(msg, fields...)
	}
	if minLevels.MinGlobalLogLevel >= INFO {
		queueAlert(INFO, msg, fields)
	}
}

func Warn(msg string, fields ...Field) {
	if minLevels.MinGlobalLogLevel >= WARN {
		zapLog.Warn(msg, fields...)
	}
	if minLevels.MinGlobalLogLevel >= WARN {
		queueAlert(WARN, msg, fields)
	}
}

func Error(msg string, fields ...Field) {
	if minLevels.MinGlobalLogLevel >= ERROR {
		zapLog.Error(msg, fields...)
	}
	if minLevels.MinGlobalLogLevel >= ERROR {
		queueAlert(ERROR, msg, fields)
	}
}

func Fatal(msg string, fields ...Field) {
	if minLevels.MinGlobalLogLevel >= FATAL {
		zapLog.Fatal(msg, fields...)
	}
	if minLevels.MinGlobalLogLevel >= FATAL {
		queueAlert(FATAL, msg, fields)
	}
}



// Set the minimum levels

func setMinimumLogLevels(minLevels *minLogLevels, logConfig *LogConfig, extLogConfig *ExternalLogConfig) {
	globalMinLevel   := util.TrimToLowerString(logConfig.GlobalMinLevel)
	externalMinLevel := util.TrimToLowerString(extLogConfig.MinLevel)

	switch globalMinLevel {
	case "", "d", "default":
		minLevels.MinGlobalLogLevel = WARN
	case "none":
		minLevels.MinGlobalLogLevel = NONE
	case "debug":
		minLevels.MinGlobalLogLevel = DEBUG
	case "info":
		minLevels.MinGlobalLogLevel = INFO
	case "warn":
		minLevels.MinGlobalLogLevel = WARN
	case "error":
		minLevels.MinGlobalLogLevel = ERROR
	case "fatal":
		minLevels.MinGlobalLogLevel = FATAL
	default:
		minLevels.MinGlobalLogLevel = WARN
		fmt.Printf("Value of minimum global log level received is not valid, min global log level is set automatically to default (WARN). [%s]", strings.TrimSpace(logConfig.GlobalMinLevel))
	}

	switch externalMinLevel {
	case "", "d", "default":
		minLevels.MinExternalLogLevel = minLevels.MinGlobalLogLevel
	case "none":
		minLevels.MinExternalLogLevel = NONE
	case "debug":
		minLevels.MinExternalLogLevel = DEBUG
	case "info":
		minLevels.MinExternalLogLevel = INFO
	case "warn":
		minLevels.MinExternalLogLevel = WARN
	case "error":
		minLevels.MinExternalLogLevel = ERROR
	case "fatal":
		minLevels.MinExternalLogLevel = FATAL
	default:
		minLevels.MinExternalLogLevel = WARN
		fmt.Printf("Value of minimum external log level received is not valid, min external log level is set automatically to default (WARN). [%s]", strings.TrimSpace(extLogConfig.MinLevel))
	}
}
