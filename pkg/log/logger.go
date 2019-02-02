package log

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger zap logger
var Logger appLogger
var once sync.Once

func init() {
	initLogger()

	// default debug level
	SetLevel(zapcore.DebugLevel)
}

// Sync flushes zap log IO
func Sync() {
	Logger.Sync()
}

// SetLevel sets the logger level
func SetLevel(l zapcore.Level) {
	Logger.atom.SetLevel(l)
}

func initLogger() {
	initLogger := func() {
		// default log level set to 'info'
		atom := zap.NewAtomicLevelAt(zap.InfoLevel)

		config := zap.Config{
			Level:       atom,
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json", // console, json, toml
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}

		mylogger, err := config.Build()
		if err != nil {
			fmt.Printf("Initialize zap logger error: %v", err)
			os.Exit(1)
		}

		Logger = appLogger{mylogger, &atom}
	}

	once.Do(initLogger)
}
