package main

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/svadikari/golang_fiber_orders/src/database"
	_ "github.com/svadikari/golang_fiber_orders/src/docs"
	"github.com/svadikari/golang_fiber_orders/src/middleware"
	orderRouters "github.com/svadikari/golang_fiber_orders/src/orders/routers"
	productRouters "github.com/svadikari/golang_fiber_orders/src/products/routers"
	"gorm.io/gorm"
)

func main() {
	// This is just a placeholder to avoid "no main function" error.
	database.ConnectDB()
	db := database.Database.Db
	if db.Error != nil {
		panic("Failed to connect to database!")
	}
	app := initApp(db)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
	slog.Info("Server is running on port 3000")
}

//@title Order, Products API
//@version 1.0
//@description This is a sample server for managing products and orders.
//@termsOfService http://swagger.io/terms/

//@contact.name API Support
//@contact.url http://www.swagger.io/support
//@contact.email shyam@shyam.com
//@license.name Apache 2.0
//@license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:3000
// @BasePath	/
func initApp(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Orders API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default status code and message
			statusCode := fiber.StatusInternalServerError
			message := "Internal Server Error"

			// Check if the error is a Fiber error
			if e, ok := err.(*fiber.Error); ok {
				statusCode = e.Code
				message = e.Message
			}

			return c.Status(fiber.StatusBadRequest).JSON(middleware.GlobalErrorHandlerResp{
				Code:   statusCode,
				Errors: strings.Split(message, ","),
			})
		},
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use(middleware.Logger)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		log := c.Locals("logger").(*slog.Logger)
		log.Info("Setting up db instance to locals")
		// Go to next middleware:
		return c.Next()
	})

	productRouters.Init(app, db)
	orderRouters.Init(app)

	return app
}
