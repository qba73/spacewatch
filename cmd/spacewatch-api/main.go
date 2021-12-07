package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/qba73/spacewatch"
)

func main() {
	if err := run(); err != nil {
		log.Printf("shutting down, error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// ========================================================================
	// Spacewatch logging setup

	log := log.New(os.Stdout, "SPACEWATCH : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// ========================================================================
	// Spacewatch configuration

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:9000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Weather struct {
			ApiKey string `conf:",noprint"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SPACEWATCH", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("SPACEWATCH", &cfg)
			if err != nil {
				return fmt.Errorf("generating config usage: %w", err)
			}
			fmt.Println(usage)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// ========================================================================
	// Spacewatch App start

	log.Println("main : Started")
	defer log.Println("main : Completed")

	cfgOut, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config: %w", err)
	}
	log.Printf("main : Config :\n%v\n", cfgOut)

	issHandler := spacewatch.ISSStatusHandler{
		ApiKey:        cfg.Weather.ApiKey,
		Log:           log,
		StatusChecker: spacewatch.GetISSStatus,
	}

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(issHandler.Get),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	// ========================================================================
	// Starting spacewatch service and listen for requests

	go func() {
		log.Printf("main : Spacewatch API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Listen for interruptions or signals from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Graceful shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("starting server: %w", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		if err != nil {
			return fmt.Errorf("could not stop spacewatch server gracefully: %w", err)
		}
	}
	return nil
}
