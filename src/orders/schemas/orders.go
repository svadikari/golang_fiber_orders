package schemas

type OrderSchema struct {
	UserId     uint              `json:"user_id" validate:"required,min=1" message:"user_id is required and must be min 1"`
	Status     OrderStatus       `json:"status" validate:"required,oneof=NEW SHIPPED DELIVERED CANCELLED"  message:"status is required and must be oneof NEW/SHIPPED/DELIVERED/CANCELLED"`
	OrderItems []OrderItemSchema `json:"order_items" validate:"required,dive" message:"order_items is required"`
}

type OrderItemSchema struct {
	ProductID uint    `json:"product_id" validate:"required,min=1" message:"product_id is required and must be min 1"`
	Quantity  int     `json:"quantity" validate:"required,min=1" message:"quantity is required and must be min 1"`
	UnitPrice float64 `json:"unit_price" validate:"required,gt=0" message:"unit_price is required and must be greater than 0"`
}

type OrderUpdateSchema struct {
	Status OrderStatus `json:"status" validate:"required,oneof=NEW SHIPPED DELIVERED CANCELLED"  message:"status is required and must be oneof NEW/SHIPPED/DELIVERED/CANCELLED"`
}

type OrderStatus string

const (
	StatusNew       OrderStatus = "NEW"
	StatusShipped   OrderStatus = "SHIPPED"
	StatusDelivered OrderStatus = "DELIVERED"
	StatusCancelled OrderStatus = "CANCELLED"
)
