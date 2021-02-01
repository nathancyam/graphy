package main

import (
	"context"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"graphy/cmd/graphy/config"
	"graphy/cmd/graphy/inject"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	v := viper.New()
	config.Init(v)
	config.Load(v)

	logger, _ := zap.NewProduction()
	srv, cleanup, err := inject.InitialiseAppServer(logger)
	if err != nil {
		logger.Fatal("failed to start application server", zap.Error(err))
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	logger.Info("server started", zap.Int("port", config.Port))

	// Cleanup when signal received.
	sig := <-done
	logger.Info("received signal to stop, draining HTTP requests", zap.String("signal", sig.String()))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cleanup()
		logger.Info("flushing logs")
		_ = logger.Sync()
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("failed to shutdown HTTP server", zap.Error(err))
	}
	logger.Info("HTTP server shutdown successfully")
}
