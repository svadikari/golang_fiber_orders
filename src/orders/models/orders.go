package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId      uint        `json:"user_id" gorm:"not null;column:user_id;index:idx_user_id"`
	TotalAmount float64     `json:"total_amount" gorm:"column:total_amount;not null;check:total_amount >= 0.1"`
	Status      string      `json:"status" gorm:"column:status;not null;size:100 enum:'NEW','SHIPPED','DELIVERED','CANCELLED' default:'NEW'"`
	OrderItems  []OrderItem `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	OrderID   uint    `json:"-" gorm:"not null;column:order_id;index:idx_order_id"`
	ProductID uint    `json:"product_id" gorm:"not null;column:product_id;index:idx_product_id"`
	Quantity  int     `json:"quantity" gorm:"column:quantity;not null;check:quantity > 0"`
	UnitPrice float64 `json:"unit_price" gorm:"column:unit_price;not null;check:unit_price >= 0.1"`
}
