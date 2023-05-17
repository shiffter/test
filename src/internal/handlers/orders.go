package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"yaa/internal/domain"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type OrdersService interface {
	GetOrder(ctx context.Context, id int64) (*domain.Order, error)
	GetOrders(ctx context.Context, o, l int) (domain.OrderSl, error)
	AddOrders(ctx context.Context, orders domain.OrderSl) error
	CompleteOrders(ctx context.Context, compOrd domain.ComplOrderSl) error
}

type Orders struct {
	service OrdersService
	logger  logrus.FieldLogger
}

func NewOrder(logger logrus.FieldLogger, service OrdersService) *Orders {
	return &Orders{
		service: service,
		logger:  logger,
	}
}

func (c *Orders) RegisterOrdersRoutes(r *mux.Router) {
	r.HandleFunc("/orders/{order_id}", c.GetOrder).Methods(http.MethodGet)
	r.HandleFunc("/orders", c.GetOrders).Methods(http.MethodGet)
	r.HandleFunc("/orders", c.AddOrders).Methods(http.MethodPost)
	r.HandleFunc("/ordcompl", c.CompleteOrders).Methods(http.MethodPost)
}

func (c *Orders) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["order_id"], 10, 64)
	if err != nil {
		c.logger.Errorf("Error converting id to int: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	order, err := c.service.GetOrder(ctx, id)
	if err != nil {
		c.logger.Errorf("Error getting order: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (c *Orders) GetOrders(w http.ResponseWriter, r *http.Request) {
	offset, limit := 0, 1
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			c.logger.Errorf("Error convert offset: %v\n", err)
			return
		}
	}

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.logger.Errorf("Error convert limit: %v\n", err)
			return
		}
	}

	ctx := r.Context()
	orders, err := c.service.GetOrders(ctx, offset, limit)
	if err != nil {
		c.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (c *Orders) AddOrders(w http.ResponseWriter, r *http.Request) {
	var OrdersSl domain.OrderSl
	err := json.NewDecoder(r.Body).Decode(&OrdersSl)
	if err != nil {
		c.logger.Error(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = c.service.AddOrders(ctx, OrdersSl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Orders) CompleteOrders(w http.ResponseWriter, r *http.Request) {
	var compOrdersSl domain.ComplOrderSl
	err := json.NewDecoder(r.Body).Decode(&compOrdersSl)
	if err != nil {
		c.logger.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = c.service.CompleteOrders(ctx, compOrdersSl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
