package domain

type Courier struct {
	Id        int64    `json:"id"`
	Type      string   `json:"type"`
	Regions   []int32  `json:"regions"`
	WorkHours []string `json:"working_hours"`
}

type Order struct {
	Id         int64    `json:"id"`
	DelivHours []string `json:"delivery_hours"`
	Cost       int32    `json:"cost"`
	Regions    int32    `json:"regions"`
	Weight     float32  `json:"weight"`
}

type CompleteOrder struct {
	IdCourier    int64  `json:"courier_id"`
	IdOrder      int64  `json:"order_id"`
	CompleteTime string `json:"completed_time"`
}

type Rating struct {
	Earn       float32
	CourRating float32
}

type CourierSl struct {
	Couriers []Courier `json:"couriers"`
}

type OrderSl struct {
	Orders []Order `json:"orders"`
}

type ComplOrderSl struct {
	CompOrd []CompleteOrder `json:"complete_orders"`
}
