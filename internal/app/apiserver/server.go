package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"wb_cource/internal/app/store"
)

type Server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	cache        *store.Cache
	orderHandler *OrderHandler
}

func newServer(st store.Store) *Server {
	cache := store.NewCache()

	s := &Server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        st,
		cache:        cache,
		orderHandler: NewOrderHandler(st, cache),
	}

	// Загружаем данные в кэш при старте
	if err := cache.LoadFromStore(st); err != nil {
		s.logger.Warn("Failed to load cache from store:", err)
	}

	s.configureRouter()
	return s
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) configureRouter() {
	// API endpoints
	s.router.HandleFunc("/order", s.orderHandler.GetOrderByID).Methods("GET")
	s.router.HandleFunc("/orders", s.orderHandler.GetAllOrders).Methods("GET")
	s.router.HandleFunc("/order", s.orderHandler.CreateOrder).Methods("POST")

	// Serve static files for web interface
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))
}
