package services

import (
	"context"
	"time"
	"yaa/internal/domain"

	"github.com/sirupsen/logrus"
)

type couriersRepo interface {
	GetCourier(ctx context.Context, id int64) (*domain.Courier, error)
	GetCouriers(ctx context.Context, o, l int) ([]domain.Courier, error)
	AddCouriers(ctx context.Context, couriers domain.CourierSl) error
	CouriersMeta(ctx context.Context, start, end time.Time, courID int64) (error, domain.Rating)
}

type CourierService struct {
	repo   couriersRepo
	logger logrus.FieldLogger
}

func NewCouriersService(repo couriersRepo, logger logrus.FieldLogger) *CourierService {
	return &CourierService{
		repo:   repo,
		logger: logger,
	}
}

func (c *CourierService) GetCourier(ctx context.Context, courierID int64) (*domain.Courier, error) {
	courier, err := c.repo.GetCourier(ctx, courierID)
	if err != nil {
		return nil, err
	}
	return courier, nil
}

func (c *CourierService) GetCouriers(ctx context.Context, o, l int) ([]domain.Courier, error) {
	couriers, err := c.repo.GetCouriers(ctx, o, l)

	if err != nil {
		return nil, err
	}
	return couriers, nil
}

func (c *CourierService) AddCouriers(ctx context.Context, couriers domain.CourierSl) error {
	err := c.repo.AddCouriers(ctx, couriers)

	if err != nil {
		return err
	}
	return nil
}

func (c *CourierService) CouriersMeta(ctx context.Context, start, end string, cour_id int64) (error, domain.Rating) {

	tStart, err := time.Parse("2006-01-02", start)
	tEnd, err := time.Parse("2006-01-02", end)

	nowStart := time.Date(tStart.Year(), tStart.Month(), tStart.Day(), 0, 0, 0, 0, time.UTC)
	nowEnd := time.Date(tEnd.Year(), tEnd.Month(), tEnd.Day(), 0, 0, 0, 0, time.UTC)
	err, result := c.repo.CouriersMeta(ctx, nowStart, nowEnd, cour_id)
	return err, result
}
