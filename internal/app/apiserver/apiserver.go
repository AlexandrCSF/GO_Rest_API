package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"wb_cource/internal/app/store"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (apiServer *APIServer) Start() error {
	if err := apiServer.configureLogger(); err != nil {
		return err
	}
	apiServer.configureRouter()
	if err := apiServer.configureStore(); err != nil {
		return err
	}
	apiServer.logger.Info("Starting API Server")

	return http.ListenAndServe(apiServer.config.BindAddr, apiServer.router)
}

func (apiServer *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(apiServer.config.LogLevel)
	if err != nil {
		return err
	}
	apiServer.logger.SetLevel(level)
	return nil
}

func (apiServer *APIServer) configureRouter() {
	apiServer.router.HandleFunc("/hello", apiServer.handleHello())
}
func (apiServer *APIServer) configureStore() error {
	st := store.New(apiServer.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	apiServer.store = st
	return nil
}
func (apiServer *APIServer) handleHello() http.HandlerFunc {
	type request struct {
		name string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World")
	}
}
