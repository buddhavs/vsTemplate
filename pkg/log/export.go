package log

import (
	"go.uber.org/zap/zapcore"
)

// Logger zap logger
var Logger appLogger

func init() {
	initLogger()
}

// Sync flushes zap log IO
func Sync() {
	Logger.Sync()
}

// SetLevel sets the logger level
func SetLevel(l zapcore.Level) {
	Logger.atom.SetLevel(l)
}
