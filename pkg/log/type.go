package log

import (
	"github.com/uber-go/zap/zapcore"
	"go.uber.org/zap"
)

type appLogger struct {
	*zap.Logger
	atom *zap.AtomicLevel
}

func (al *appLogger) SetLevel(l zapcore.Level) {
	al.atom.SetLevel(l)
}
