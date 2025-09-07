package services

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/svadikari/golang_fiber_orders/src/products/models"
)

type mockProductRepository struct {
	mock.Mock
}

func (r *mockProductRepository) Create(product *models.Product) *models.Product {
	args := r.Called(product)
	return args.Get(0).(*models.Product)
}

func (m *mockProductRepository) Update(product *models.Product) *models.Product {
	args := m.Called(product)
	return args.Get(0).(*models.Product)
}

func (m *mockProductRepository) Delete(product *models.Product) {
	m.Called(product)
}

func (m *mockProductRepository) Find() []models.Product {
	args := m.Called()
	return args.Get(0).([]models.Product)
}

func (m *mockProductRepository) FindByID(id uint) models.Product {
	args := m.Called(id)
	return args.Get(0).(models.Product)
}

func TestGetProductById(t *testing.T) {

	mockRepo := new(mockProductRepository)
	service := NewProductService(mockRepo, slog.Default())

	t.Run("Product doesn't found for given productId", func(t *testing.T) {
		product := models.Product{}
		mockRepo.On("FindByID", uint(1)).Return(product).Once()
		_, err := service.GetProductByID(1)
		assert.Error(t, err)
		assert.Equal(t, "product not found for Id:1", err.Error())
		mockRepo.AssertExpectations(t)
		mockRepo.On("FindByID", uint(1)).Unset()
	})

	t.Run("Product Found for given productId", func(t *testing.T) {
		product := models.Product{Name: "Test Product", Price: 10.0, Stock: 100, Description: "Test Description", ImageURL: "http://example.com/image.jpg"}
		product.ID = 1
		mockRepo.On("FindByID", uint(1)).Return(product).Once()
		result, err := service.GetProductByID(1)
		assert.NoError(t, err)
		assert.Equal(t, product, result)
		mockRepo.AssertExpectations(t)
		mockRepo.On("FindByID", uint(1)).Unset()
	})

}
