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
}

func NewOrderHandler(store store.Store, cache *store.Cache) *OrderHandler {
	return &OrderHandler{
		store: store,
		cache: cache,
	}
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("order_uid")
	if orderUID == "" {
		http.Error(w, "order_uid parameter is required", http.StatusBadRequest)
		return
	}

	// Сначала проверяем кэш
	if order, exists := h.cache.Get(orderUID); exists {
		respondWithJSON(w, http.StatusOK, order)
		return
	}

	// Если в кэше нет, ищем в БД
	order, err := h.store.Order().FindByID(orderUID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Сохраняем в кэш
	h.cache.Set(order)

	respondWithJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders := h.cache.GetAll()
	respondWithJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.store.Order().Create(&order); err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Сохраняем в кэш
	h.cache.Set(&order)

	respondWithJSON(w, http.StatusCreated, order)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
