package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/svadikari/golang_fiber_orders/src/orders/controllers"
)

func Init(app *fiber.App) {
	app.Route("/orders", func(router fiber.Router) {
		router.Get("/", controllers.GetOrders)
		router.Post("/", controllers.CreateOrders)
		router.Put("/:id", controllers.UpdateOrder)
		router.Get("/:id", controllers.GetOrder)
		router.Delete("/:id", controllers.DeleteOrder)
		router.Get("/consumer/start", controllers.StartConsumer)
		router.Get("/consumer/stop", controllers.StopConsumer)
	})
}
