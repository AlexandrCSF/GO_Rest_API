package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"wb_cource/internal/app/store"
)

type Server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *Server {
	s := &Server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}
	s.configureRouter()
	return s
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) configureRouter() {

}
