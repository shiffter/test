package services

import (
	"context"
	"yaa/internal/domain"

	"github.com/sirupsen/logrus"
)

type ordersRepo interface {
	GetOrder(ctx context.Context, id int64) (*domain.Order, error)
	GetOrders(ctx context.Context, o, l int) (domain.OrderSl, error)
	AddOrders(ctx context.Context, orders domain.OrderSl) error
	ExistOrder(ctx context.Context, c, o int64) (bool, error)
	SetCompliteOrders(ctx context.Context, o, c int64, str string) error
}

type OrderService struct {
	repo   ordersRepo
	logger logrus.FieldLogger
}

func NewOrderService(repo ordersRepo, logger logrus.FieldLogger) *OrderService {
	return &OrderService{
		repo:   repo,
		logger: logger,
	}
}

func (c *OrderService) GetOrder(ctx context.Context, orderID int64) (*domain.Order, error) {
	order, err := c.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (c *OrderService) GetOrders(ctx context.Context, o, l int) (domain.OrderSl, error) {
	orders, err := c.repo.GetOrders(ctx, o, l)

	if err != nil {
		return domain.OrderSl{}, err
	}
	return orders, nil
}

func (c *OrderService) AddOrders(ctx context.Context, orders domain.OrderSl) error {
	err := c.repo.AddOrders(ctx, orders)

	if err != nil {
		return err
	}
	return nil
}

func (c *OrderService) CompleteOrders(ctx context.Context, ord domain.ComplOrderSl) error {
	for _, order := range ord.CompOrd {
		vars, err := c.repo.ExistOrder(ctx, order.IdCourier, order.IdOrder)
		if vars && err == nil {
			err = c.repo.SetCompliteOrders(ctx, order.IdCourier, order.IdOrder, order.CompleteTime)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
