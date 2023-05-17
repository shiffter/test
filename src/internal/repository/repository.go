package repository

import (
	"context"
	"time"
	"yaa/internal/domain"
	"yaa/internal/repository/queries"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	GetCourier(ctx context.Context, id int64) (*domain.Courier, error)
	GetCouriers(ctx context.Context, o, l int) ([]domain.Courier, error)
	AddCouriers(ctx context.Context, couriers domain.CourierSl) error
	GetOrder(ctx context.Context, id int64) (*domain.Order, error)
	GetOrders(ctx context.Context, o, l int) (domain.OrderSl, error)
	AddOrders(ctx context.Context, orders domain.OrderSl) error
	ExistOrder(ctx context.Context, c, o int64) (bool, error)
	SetCompliteOrders(ctx context.Context, c, o int64, str string) error
	CouriersMeta(ctx context.Context, start, end time.Time, courID int64) (error, domain.Rating)
}

type repo struct {
	*queries.Queries
	logger logrus.FieldLogger
	pool   pgxpool.Pool
}

func NewRepository(pgxPool *pgxpool.Pool, logger logrus.FieldLogger) Repository {
	return &repo{
		Queries: queries.New(pgxPool),
		logger:  logger,
		pool:    *pgxPool,
	}
}
