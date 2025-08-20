package apiserver

import (
	"net/http"
	"wb_cource/internal/app/store/sqlstore"
)

func Start(config *Config) error {
	store := sqlstore.New(config.DatabaseURL)
	if err := store.Open(); err != nil {
		return err
	}
	defer store.Close()

	srv := newServer(store)
	return http.ListenAndServe(config.BindAddr, srv)
}
