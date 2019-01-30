package log

import (
	"go.uber.org/zap"
)

type appLogger struct {
	*zap.Logger
	atom *zap.AtomicLevel
}
