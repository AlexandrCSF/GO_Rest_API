package apiserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"wb_cource/internal/app/config"
	"wb_cource/internal/app/store/sqlstore"
)

func Start(cfg *config.Config) error {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithField("addr", cfg.BindAddr).Info("server_starting")

	store := sqlstore.New(cfg.DatabaseURL)
	if err := store.Open(); err != nil {
		logger.WithError(err).Error("db_connect_failed")
		return err
	}
	logger.Info("db_connected")
	defer store.Close()

	srv := newServer(store, cfg)

	httpServer := &http.Server{
		Addr:    cfg.BindAddr,
		Handler: srv,
	}

	srv.StartKafkaConsumer()

	go func() {
		logger.WithField("addr", cfg.BindAddr).Info("server_listen")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Error("server_error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("server_shutting_down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.Stop()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("server_shutdown_error")
		return err
	}

	logger.Info("server_exited")
	return nil
}
