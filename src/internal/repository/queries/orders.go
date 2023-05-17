package queries

import (
	"context"
	"fmt"
	"time"
	"yaa/internal/domain"
)

func (r *Queries) GetOrder(ctx context.Context, id int64) (*domain.Order, error) {
	query := "SELECT id, delivery_hours, cost, regions, weight FROM orders where id = $1"
	rows := r.pool.QueryRow(ctx, query, id)

	var c domain.Order

	err := rows.Scan(&c.Id, &c.DelivHours, &c.Cost, &c.Regions, &c.Weight)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *Queries) GetOrders(ctx context.Context, offset, limit int) (domain.OrderSl, error) {
	query := "SELECT id, delivery_hours, cost, regions, weight FROM orders ORDER BY id OFFSET $1 LIMIT $2"
	rows, err := r.pool.Query(ctx, query, offset, limit)
	if err != nil {
		return domain.OrderSl{}, err
	}

	defer rows.Close()

	var cours domain.OrderSl

	for rows.Next() {
		var c domain.Order
		err = rows.Scan(&c.Id, &c.DelivHours, &c.Cost, &c.Regions, &c.Weight)
		if err != nil {
			return domain.OrderSl{}, err
		}

		cours.Orders = append(cours.Orders, c)
	}

	return cours, nil
}

func (r *Queries) AddOrders(ctx context.Context, orders domain.OrderSl) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	stmt, err := tx.Prepare(ctx, "insert_ord", `INSERT INTO orders (id, delivery_hours, cost, regions, weight, completed_time) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}
	for _, v := range orders.Orders {
		_, err := tx.Exec(ctx, stmt.SQL, v.Id, v.DelivHours, v.Cost, v.Regions, v.Weight, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Queries) CheckOrderStatus(ctx context.Context, ordID int64) (bool, error) {
	var t time.Time

	rows, err := r.pool.Query(context.Background(), "SELECT completed_time from orders where id=$1", ordID)
	defer rows.Close()
	if err != nil {
		return false, fmt.Errorf("err")
	}
	rows.Next()
	rows.Scan(&t)
	if !t.IsZero() {
		err = fmt.Errorf("order already done")
		return false, err
	}
	return true, nil
}

func (r *Queries) ExistOrder(ctx context.Context, courID, orderID int64) (bool, error) {
	status, err := r.CheckOrderStatus(ctx, orderID)

	if !status {
		err = fmt.Errorf("order have complete time")
		return false, err
	}
	if err != nil {
		return false, err
	}

	var existsOrd bool
	err = r.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM orders WHERE id =$1 )", orderID).Scan(&existsOrd)
	if err != nil {
		return false, err
	}

	var existsCour bool
	err = r.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM couriers WHERE id =$1 )", courID).Scan(&existsCour)
	if existsCour && existsOrd {
		return true, nil
	}
	return false, nil

}

func (r *Queries) SetCompliteOrders(ctx context.Context, c, o int64, str string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	stmt, err := tx.Prepare(ctx, "prod", `INSERT INTO complete_orders (courier_id, order_id,  completed_time) VALUES ($1, $2, $3)`)
	if err != nil {
		return err
	}
	now := time.Now()
	t, err := time.Parse("15:04", str)
	year, month, _ := now.Date()
	now = time.Date(year, month, 1, t.Hour(), t.Minute(), 0, 0, time.UTC)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, stmt.SQL, c, o, now)

	if err != nil {
		return err
	}

	stmtStatus, err := tx.Prepare(ctx, "prod1", `UPDATE orders SET completed_time = ($1) where id = ($2)`)

	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, stmtStatus.SQL, now, o)

	if err != nil {
		return err
	}

	return err
}
