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

type CouriersService interface {
	GetCourier(ctx context.Context, id int64) (*domain.Courier, error)
	GetCouriers(ctx context.Context, o, l int) ([]domain.Courier, error)
	AddCouriers(ctx context.Context, couriers domain.CourierSl) error
	CouriersMeta(ctx context.Context, start, end string, courID int64) (error, domain.Rating)
}

type Couriers struct {
	service CouriersService
	logger  logrus.FieldLogger
}

func NewCourier(logger logrus.FieldLogger, service CouriersService) *Couriers {
	return &Couriers{
		service: service,
		logger:  logger,
	}
}

func (c *Couriers) RegisterCouriersRoutes(r *mux.Router) {
	r.HandleFunc("/couriers/{courier_id}", c.GetCourier).Methods(http.MethodGet)
	r.HandleFunc("/couriers", c.GetCouriers).Methods(http.MethodGet)
	r.HandleFunc("/couriers", c.AddCouriers).Methods(http.MethodPost)
	r.HandleFunc("/couriers/meta-info/{courier_id}", c.CouriersMeta).Methods(http.MethodGet)
}

func (c *Couriers) GetCourier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["courier_id"], 10, 64)
	if err != nil {
		c.logger.Error("Error converting id to int: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := c.service.GetCourier(ctx, id)
	if err != nil {
		c.logger.Error("Error getting user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (c *Couriers) GetCouriers(w http.ResponseWriter, r *http.Request) {
	offset, limit := 0, 1
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			c.logger.Error("Error convert offset: %v\n", err)
			return
		}
	}

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.logger.Error("Error convert limit: %v\n", err)
			return
		}
	}

	ctx := r.Context()
	couriers, err := c.service.GetCouriers(ctx, offset, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(couriers)
}

func (c *Couriers) AddCouriers(w http.ResponseWriter, r *http.Request) {
	var CourSl domain.CourierSl
	err := json.NewDecoder(r.Body).Decode(&CourSl)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = c.service.AddCouriers(ctx, CourSl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Couriers) CouriersMeta(w http.ResponseWriter, r *http.Request) {
	var start, end string
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["courier_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	start = r.URL.Query().Get("start_date")
	end = r.URL.Query().Get("end_date")
	ctx := r.Context()
	err, result := c.service.CouriersMeta(ctx, start, end, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
