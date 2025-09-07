package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null;column:name;index:idx_name;size:200"`
	Description string  `json:"description" gorm:"column:description;not null;size:1000"`
	Price       float64 `json:"price" gorm:"column:price;not null;check:price >= 0.1"`
	Stock       int     `json:"stock" gorm:"column:stock"`
	ImageURL    string  `json:"image_url" gorm:"column:image_url"`
}
