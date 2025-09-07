package repository

import (
	"errors"
	"log"

	"github.com/svadikari/golang_fiber_orders/src/products/models"
	"gorm.io/gorm"
)

type productRepository struct {
	Db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{Db: db}
}

func (r *productRepository) Create(product *models.Product) *models.Product {
	r.Db.Create(&product)
	return product
}

func (r *productRepository) Find() []models.Product {
	var products []models.Product

	r.Db.Find(&products) // Replace nil with actual model slice
	return products
}

func (r *productRepository) FindByID(id uint) models.Product {
	var product models.Product
	result := r.Db.First(&product, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Product{}
	}
	return product
}

func (r *productRepository) Update(product *models.Product) *models.Product {
	r.Db.Save(&product)
	return product
}

func (r *productRepository) Delete(product *models.Product) {
	result := r.Db.Delete(&product)
	log.Println(result.Error.Error())
}

type ProductRepository interface {
	Find() []models.Product
	FindByID(uint) models.Product
	Create(*models.Product) *models.Product
	Update(*models.Product) *models.Product
	Delete(*models.Product)
}
