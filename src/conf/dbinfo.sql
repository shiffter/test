create type courier_type as ENUM ('FOOT', 'BIKE', 'AUTO');

create table if not exists couriers (
	id BIGINT PRIMARY KEY,
	cour_type courier_type NOT NULL,
	regions int4[],
	working_hours TEXT[]
);

create table if not exists orders (
	id BIGINT PRIMARY KEY UNIQUE,
	delivery_hours TEXT[],
	cost int,
	regions int,
	weight float,
	completed_time timestamp
);


CREATE TABLE if not exists complete_orders (
    courier_id BIGINT NOT NULL REFERENCES couriers(id),
    order_id  BIGINT NOT NULL REFERENCES orders(id) UNIQUE,
    completed_time TIMESTAMP NOT NULL
);
