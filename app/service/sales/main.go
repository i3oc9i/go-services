package main

import (
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ardanlabs/conf"
	"github.com/i3oc9i/go-service/app/service/sales/handlers"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/* TODO
- Need to figure out timeouts for the http service.
*/

var build = "develop"

func main() {

	log, err := newLogger("SALES")
	if err != nil {
		fmt.Println("Fatal error constructing logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := setup(log); err != nil {
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

func setup(log *zap.SugaredLogger) error {

	// -------------------------------------------------------------- GOMAXPROCS
	// Set the correct number of threads based on how many CPUs are available
	// either by the machine or quotas.

	opt := maxprocs.Logger(log.Infof)

	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------- Configuration
	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string `conf:"default:0.0.0.0:3000"`
			DebugHost       string `conf:"default:0.0.0.0:4000"`
			ReadTimeout     string `conf:"default:5s"`
			WriteTimeout    string `conf:"default:10s"`
			IdleTimeout     string `conf:"default:120s"`
			ShutdownTimeout string `conf:"default:20s"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "Sales services",
		},
	}

	if help, err := conf.ParseOSArgs("SALES", &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------- App Starting
	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for the output: %w", err)
	}
	log.Infow("startup", "config", out)

	expvar.NewString("build").Set(build)

	// -------------------------------------------------------------- Start Debug Service
	log.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	// The Debug function returns a mux to listen and serve on for all the debug
	// related endpoints. This includes the standard library endpoints.
	debugMux := handlers.DebugStandardLibraryMux()

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// -------------------------------------------------------------- Gracefully Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	return nil

}
