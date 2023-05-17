package queries

import (
	"context"
	"time"
	"yaa/internal/domain"

	"github.com/jackc/pgtype"
)

func (r *Queries) GetCourier(ctx context.Context, id int64) (*domain.Courier, error) {
	query := "SELECT id, cour_type, regions, working_hours FROM couriers where id = $1"
	rows := r.pool.QueryRow(ctx, query, id)
	var c domain.Courier
	err := rows.Scan(&c.Id, &c.Type, &c.Regions, &c.WorkHours)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Queries) GetCouriers(ctx context.Context, offset, limit int) ([]domain.Courier, error) {
	query := "SELECT id, cour_type, regions, working_hours FROM couriers ORDER BY id OFFSET $1 LIMIT $2"
	rows, err := r.pool.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cours []domain.Courier

	for rows.Next() {
		var c domain.Courier
		err = rows.Scan(&c.Id, &c.Type, &c.Regions, &c.WorkHours)
		if err != nil {
			return nil, err
		}
		cours = append(cours, c)
	}
	return cours, nil
}

func (r *Queries) AddCouriers(ctx context.Context, couriers domain.CourierSl) error {
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

	stmt, err := tx.Prepare(ctx, "prod", `INSERT INTO couriers (id, cour_type, regions,
	working_hours) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	for _, v := range couriers.Couriers {
		_, err := tx.Exec(ctx, stmt.SQL, v.Id, v.Type, v.Regions, v.WorkHours)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Queries) CouriersMeta(ctx context.Context, start, end time.Time, courID int64) (error, domain.Rating) {
	res := domain.Rating{}
	err, earn := r.getEarn(ctx, start, end, courID)
	if err != nil {
		return err, res
	}
	err, rating := r.getRating(ctx, start, end, courID)
	if err != nil {
		return err, res
	}
	res.Earn, res.CourRating = earn, rating
	return err, res
}

func (r *Queries) getRating(ctx context.Context, start, end time.Time, courID int64) (error, float32) {
	query := `with cour_type as (
	select id, cour_type,
	case 
		when c.cour_type = 'AUTO' then 1.0
		when c.cour_type = 'BIKE' then 2.0
		when c.cour_type = 'FOOT' then 3.0
		else 0
		end as type_category
	from couriers c where id = $3
)
SELECT ((COUNT(*) / EXTRACT(EPOCH FROM (CAST($2 AS timestamp) - CAST($1 AS timestamp))/3600)) * cour_type.type_category)::float4 AS orders_per_hour
FROM complete_orders
JOIN cour_type ON complete_orders.courier_id = cour_type.id
WHERE courier_id = $3 AND completed_time >= $1 AND completed_time < $2
GROUP BY cour_type.type_category;`

	row := r.pool.QueryRow(ctx, query, start, end, courID)
	var result pgtype.Float4
	err := row.Scan(&result)
	if err != nil {
		return err, 0.0
	}
	if result.Status == pgtype.Null {
		return nil, 0.0
	}
	return err, result.Float
}
func (r *Queries) getEarn(ctx context.Context, start, end time.Time, courID int64) (error, float32) {

	query := `
		with cour_type as (
					select id, cour_type,              
					case 
									when c.cour_type = 'AUTO' then 4.0
									when c.cour_type = 'BIKE' then 3.0
									when c.cour_type = 'FOOT' then 2.0
									else 0
									end as type_category
					from couriers c where id = $3
	)         

	SELECT SUM(o.cost * cour_type.type_category)::float4 as total_cost
	FROM complete_orders co
	JOIN orders o ON co.order_id = o.id
	JOIN cour_type ON co.courier_id = cour_type.id
	WHERE co.courier_id = $3 AND co.completed_time >= $1 AND co.completed_time < $2;`

	row := r.pool.QueryRow(ctx, query, start, end, courID)
	var result pgtype.Float4
	err := row.Scan(&result)

	if err != nil {
		return err, 0.0
	}
	if result.Status == pgtype.Null {
		return nil, 0.0
	}
	return err, result.Float
}
