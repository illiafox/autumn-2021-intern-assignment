package zap

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(log io.Writer) *zap.Logger {
	pe := zap.NewProductionEncoderConfig()

	pe.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC1123)
	fileEncoder := zapcore.NewJSONEncoder(pe)

	pe.EncodeTime = zapcore.TimeEncoderOfLayout("02/01/2006 15:04:05")
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(log), zap.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)

	return zap.New(core)
}
