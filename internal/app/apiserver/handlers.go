package apiserver

import (
	"encoding/json"
	"net/http"
	"wb_cource/internal/app/model"
	"wb_cource/internal/app/store"
)

type OrderHandler struct {
	store store.Store
	cache *store.Cache
	log   *ServerLogger
}

// ServerLogger — тонкая обёртка, чтобы не тянуть logrus напрямую в хэндлер
type ServerLogger struct{ s *Server }

func (l *ServerLogger) Error(msg string) { l.s.logger.Error(msg) }
func (l *ServerLogger) Warn(msg string)  { l.s.logger.Warn(msg) }
func (l *ServerLogger) Info(msg string)  { l.s.logger.Info(msg) }

func NewOrderHandler(store store.Store, cache *store.Cache) *OrderHandler {
	return &OrderHandler{
		store: store,
		cache: cache,
	}
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("order_uid")
	if orderUID == "" {
		h.respondError(w, http.StatusBadRequest, "order_uid parameter is required")
		return
	}

	if order, exists := h.cache.Get(orderUID); exists {
		h.respondJSON(w, http.StatusOK, order)
		return
	}

	order, err := h.store.Order().FindByID(orderUID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "order not found")
		return
	}

	h.cache.Set(order)
	h.respondJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders := h.cache.GetAll()
	h.respondJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.store.Order().Create(&order); err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	h.cache.Set(&order)
	h.respondJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *OrderHandler) respondError(w http.ResponseWriter, code int, msg string) {
	h.respondJSON(w, code, map[string]string{"error": msg})
}
