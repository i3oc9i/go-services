package main

import (
	"fmt"
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.uber.org/automaxprocs/maxprocs"
)

var build = "develop"

func main() {

	log, err := newLogger("SALES")
	if err != nil {
		fmt.Println("Fatal error constructing logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func newLogger(service string) (*zap.SugaredLogger, error) {

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}

func run(log *zap.SugaredLogger) error {

	if _, err := maxprocs.Set(0); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS())

	return nil
}
