package apiserver

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	cache := store.NewCache()

	s := &Server{
		router:       mux.NewRouter(),
		logger:       logger,
		store:        st,
		cache:        cache,
		orderHandler: NewOrderHandler(st, cache),
	}

	if err := cache.LoadFromStore(st); err != nil {
		s.logger.WithError(err).Warn("cache_warmup_failed")
	}

	s.configureRouter()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) configureRouter() {
	s.router.Use(s.accessLogMiddleware)

	s.router.HandleFunc("/order", s.orderHandler.GetOrderByID).Methods("GET")
	s.router.HandleFunc("/orders", s.orderHandler.GetAllOrders).Methods("GET")
	s.router.HandleFunc("/order", s.orderHandler.CreateOrder).Methods("POST")

	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (s *Server) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		dur := time.Since(start)

		s.logger.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      rec.status,
			"duration_ms": dur.Milliseconds(),
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		}).Info("http_request")
	})
}
