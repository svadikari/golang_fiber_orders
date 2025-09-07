package controllers

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/svadikari/golang_fiber_orders/src/products/schemas"
	"github.com/svadikari/golang_fiber_orders/src/products/services"
)

type ProductController interface {
	GetProducts(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	GetProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type productController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) ProductController {
	return &productController{productService: productService}
}

// Get all products
//
//	@Summary		Get all products
//	@Description	Retrieve a list of all products
//	@Tags			Products
//	@Produce		json
//	@Success		200	{array}	schemas.Product
//	@Router			/products [get]
func (pc *productController) GetProducts(c *fiber.Ctx) error {
	products, _ := pc.productService.GetAllProducts()
	return c.Status(fiber.StatusOK).JSON(products)
}

// Create product
//
//	@Summary		Create product
//	@Description	Creates a new product
//	@Tags			Products
//
//	@Accept			json
//
//	@Produce		json
//
//	@Param			product	body		schemas.Product	true	"Product payload"
//
//	@Success		201		{object}	schemas.Product
//	@Router			/products [post]
func (pc *productController) CreateProduct(c *fiber.Ctx) error {
	var productPayload schemas.Product
	log := c.Locals("logger").(*slog.Logger)
	if err := c.BodyParser(&productPayload); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"details": err.Error(),
		})
	}
	product, err := pc.productService.CreateProduct(productPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(product)
}

// Get product
//
//	@Summary		Get product
//	@Description	Get a product by ID
//	@Tags			Products
//
//	@Accept			json
//
//	@Produce		json
//
//	@Param			id	path		int	true	"Product ID"
//
//	@Success		200	{object}	schemas.Product
//	@Router			/products/{id} [get]
func (pc *productController) GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"details": fmt.Sprintf("Invalid product ID parameter: %v", err.Error()),
		})
	}
	log := c.Locals("logger").(*slog.Logger)
	log.Info("Fetching product by ID", "id", id)
	product, err := pc.productService.GetProductByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(product)
}

// Update product
//
//	@Summary		Update product
//	@Description	Update a new product
//	@Tags			Products
//
//	@Accept			json
//
//	@Produce		json
//
//	@Param			product	body		schemas.Product	true	"Product payload"
//	@Param			id		path		int				true	"Product ID"
//
//	@Success		200		{object}	schemas.Product
//	@Router			/products/{id} [put]
func (pc *productController) UpdateProduct(c *fiber.Ctx) error {
	log := c.Locals("logger").(*slog.Logger)

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error("Invalid product ID parameter", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"details": fmt.Sprintf("Invalid product ID parameter: %v", err.Error()),
		})
	}
	var productPayload schemas.Product

	if err := c.BodyParser(&productPayload); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"details": err.Error(),
		})
	}
	product, err := pc.productService.UpdateProduct(uint(id), productPayload)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// Delete product
//
//	@Summary		Delete product
//	@Description	Delete product by ID
//	@Tags			Products
//
//	@Accept			json
//
//	@Produce		json
//
//	@Param			id	path	int	true	"Product ID"
//
//	@Success		204
//	@Router			/products/{id} [delete]
func (pc *productController) DeleteProduct(c *fiber.Ctx) error {
	log := c.Locals("logger").(*slog.Logger)

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error("Invalid product ID parameter", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"details": fmt.Sprintf("Invalid product ID parameter: %v", err.Error()),
		})
	}
	err = pc.productService.DeleteProduct(uint(id))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"details": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
