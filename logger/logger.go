package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"logfile.log"}

	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	logger, err := cfg.Build(zap.AddCaller())
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	return logger
}
