package apiserver

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"wb_cource/internal/app/store/sqlstore"
)

func Start(config *Config) error {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithField("addr", config.BindAddr).Info("server_starting")

	store := sqlstore.New(config.DatabaseURL)
	if err := store.Open(); err != nil {
		logger.WithError(err).Error("db_connect_failed")
		return err
	}
	logger.Info("db_connected")
	defer store.Close()

	srv := newServer(store)
	logger.WithField("addr", config.BindAddr).Info("server_listen")
	return http.ListenAndServe(config.BindAddr, srv)
}
