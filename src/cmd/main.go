package main

import (
	"log"
	"net/http"
	"os"
	"yaa/internal/handlers"
	"yaa/internal/repository"
	"yaa/internal/services"
	"yaa/pkg/postgres"

	"golang.org/x/time/rate"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {

	logger := logrus.New()

	pool, err := postgres.NewPool(os.Getenv("POSTGRES_DSN"))

	if err != nil {
		logger.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewRepository(pool, logger)
	courierService := services.NewCouriersService(repo, logger)
	orderService := services.NewOrderService(repo, logger)

	r := mux.NewRouter()

	limiter := rate.NewLimiter(1, 10)
	r.Use(rateLimitMiddleware(limiter))

	courierHandler := handlers.NewCourier(logger, courierService)
	courierHandler.RegisterCouriersRoutes(r)

	orderHandler := handlers.NewOrder(logger, orderService)
	orderHandler.RegisterOrdersRoutes(r)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func rateLimitMiddleware(limiter *rate.Limiter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter.Allow() == false {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
