package services

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/svadikari/golang_fiber_orders/src/products/models"
	"github.com/svadikari/golang_fiber_orders/src/products/repository"
	"github.com/svadikari/golang_fiber_orders/src/products/schemas"
	"gorm.io/gorm"
)

type ProductService interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByID(uint) (models.Product, error)
	CreateProduct(schemas.Product) (models.Product, error)
	UpdateProduct(uint, schemas.Product) (models.Product, error)
	DeleteProduct(uint) error
}

type productService struct {
	Db                *gorm.DB
	Logger            *slog.Logger
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository, logger *slog.Logger) ProductService {
	logger = logger.With("service", "ProductService")
	return &productService{Db: nil, Logger: logger, productRepository: productRepository}
}

func (s *productService) CreateProduct(productPaylod schemas.Product) (models.Product, error) {
	product := models.Product{
		Name:        productPaylod.Name,
		Description: productPaylod.Description,
		Price:       productPaylod.Price,
		Stock:       productPaylod.Stock,
		ImageURL:    productPaylod.ImageURL,
	}
	s.productRepository.Create(&product)
	s.Logger.Info("Created new product in the database", "product", product)
	return product, nil
}

func (s *productService) GetAllProducts() ([]models.Product, error) {
	s.Logger.Info("Fetching all products from the database")
	products := s.productRepository.Find()
	s.Logger.Info("Fetched products from the database", "products", products)
	return products, nil
}

func (s *productService) GetProductByID(id uint) (models.Product, error) {
	s.Logger.Info("Fetching product by ID from the database", "id", id)
	product := s.productRepository.FindByID(id)
	if product.ID == 0 {
		return product, fmt.Errorf("product not found for Id:%v", id)
	}
	return product, nil
}

func (s *productService) UpdateProduct(id uint, productPaylod schemas.Product) (models.Product, error) {
	product := s.productRepository.FindByID(id)
	if product.ID == 0 {
		s.Logger.Warn("Product not found in the database", "id", id)
		return models.Product{}, errors.New("Product not found for ID: " + strconv.Itoa(int(id)))
	}

	if productPaylod.Name != "" {
		product.Name = productPaylod.Name
	}
	if productPaylod.Description != "" {
		product.Description = productPaylod.Description
	}
	if productPaylod.Price != 0 {
		product.Price = productPaylod.Price
	}
	if productPaylod.Stock != 0 {
		product.Stock = productPaylod.Stock
	}
	if productPaylod.ImageURL != "" {
		product.ImageURL = productPaylod.ImageURL
	}
	product.ID = id
	s.productRepository.Update(&product)
	s.Logger.Info("Created new product in the database", "product", product)
	return product, nil
}

func (s *productService) DeleteProduct(id uint) error {
	product := s.productRepository.FindByID(id)
	if product.ID == 0 {
		s.Logger.Warn("Product not found in the database", "id", id)
		return errors.New("Product not found for ID: " + strconv.Itoa(int(id)))
	}
	s.productRepository.Delete(&product)
	return nil
}
