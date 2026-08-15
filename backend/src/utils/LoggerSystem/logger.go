package logger_system


import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func LoggerInit(env string) {
	ConsoleEncoderCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		EncodeTime:     zapcore.TimeEncoderOfLayout("[2006/01/02] [15:04:05]"),
		EncodeLevel:    zapcore.CapitalLevelEncoder, 
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		ConsoleSeparator: " ",
	}
	
	JSONEncoderCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		EncodeTime:     zapcore.TimeEncoderOfLayout("[2006/01/02] [15:04:05]"),
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		ConsoleSeparator: " ",
	}

	cwd, _ := os.Getwd()
	fp := fmt.Sprintf("%s/logs", cwd)

	os.MkdirAll(fp, 0755)

	rotation := &lumberjack.Logger{
		Filename: "logs/app.log",
		MaxSize: 10,
		MaxBackups: 2,
		MaxAge: 30,
		Compress: true,
	}

	var (
		consoleCore zapcore.Core
		fileCore zapcore.Core
	)

	if env == "production" {
		ConsoleEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleCore = zapcore.NewCore(
			zapcore.NewConsoleEncoder(ConsoleEncoderCfg),
			zapcore.AddSync(os.Stderr),
			zap.InfoLevel,
		)
	} else {
		ConsoleEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleCore = zapcore.NewCore(
			zapcore.NewConsoleEncoder(ConsoleEncoderCfg),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		)
	}

	fileCore = zapcore.NewCore(
		zapcore.NewJSONEncoder(JSONEncoderCfg),  
		zapcore.AddSync(rotation),             
		zap.InfoLevel,
	)

	Log = zap.New(zapcore.NewTee(consoleCore, fileCore),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(zap.Int("pid", os.Getpid())),
	)
}