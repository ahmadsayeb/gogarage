package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/service/foundation/logger"
)

var build = "develop"

// pattern used here is to have the main function that calls
// a run function. if error, log error and exit
func main() {
	var log *logger.Logger
	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* Send Alert ********")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return "" //web.GetTraceID(ctx)
	}
	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events)
	//-------
	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// GOMAXPROX
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:5s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:15s"`
			APIHost            string        `conf:"default:0.0.0.:3000"`
			DebugHost          string        `conf:"default:0.0.0.:4000"`
			CORSAllowedOrigins []string      `conf:"default:*,mask"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Sales",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		// sentiner error: don't continue the prorgram, user asked for help
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	//--------------------------------------------------------------------------------------------------------------
	// App STARTING
	log.Info(ctx, "starting services", "version", cfg.Build) // The struct is embedded so you don't need to use cfg.Version.Build
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)
	//--------------------------------------------------------------------------------------------------------------
	// APP STARTING

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdown
	log.Info(ctx, "suthdown", "status", "shutdown started", "signal", sig)
	defer log.Info(ctx, "suthdown", "status", "shutdown started", "signal", sig)

	return nil
}
