package log

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
)

var once sync.Once

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
