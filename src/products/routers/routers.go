package routers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/svadikari/golang_fiber_orders/src/products/controllers"
	"github.com/svadikari/golang_fiber_orders/src/products/repository"
	"github.com/svadikari/golang_fiber_orders/src/products/services"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB) {
	productController := initializeFramework(db)
	app.Route("/products", func(router fiber.Router) {
		router.Get("/", productController.GetProducts)
		router.Post("/", productController.CreateProduct)
		router.Put("/:id<min(1)>", productController.UpdateProduct)
		router.Get("/:id<min(1)>", productController.GetProduct)
		router.Delete("/:id", productController.DeleteProduct)
	})
}

func initializeFramework(db *gorm.DB) controllers.ProductController {
	productRepository := repository.NewProductRepository(db)
	productService := services.NewProductService(productRepository, slog.Default())
	return controllers.NewProductController(productService)
}
