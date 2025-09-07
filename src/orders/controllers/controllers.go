package controllers

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/svadikari/golang_fiber_orders/src/middleware"
	"github.com/svadikari/golang_fiber_orders/src/orders/models"
	"github.com/svadikari/golang_fiber_orders/src/orders/schemas"
	"gorm.io/gorm"
)

// Get all orders
//
//	@Summary		Get all orders
//	@Description	Retrieve a list of all orders
//	@Tags			Orders
//	@Produce		json
//	@Success		200	{array}	models.Order
//	@Router			/orders [get]
func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order
	db := c.Locals("db").(*gorm.DB)
	db.Preload("OrderItems").Find(&orders)
	return c.Status(fiber.StatusOK).JSON(orders)
}

// Create Order
//
//	@Summary		Create Order
//	@Description	Creates a new order
//	@Tags			Orders
//	@Produce		json
//	@Accept			json
//	@Param			order	body		schemas.OrderSchema	true	"Order payload"
//	@Success		200		{object}	models.Order
//
//	@Failure		400		{object}	middleware.GlobalErrorHandlerResp
//	@Failure		404		{object}	middleware.GlobalErrorHandlerResp
//	@Failure		500		{object}	middleware.GlobalErrorHandlerResp
//
//	@Router			/orders [post]
func CreateOrders(c *fiber.Ctx) error {
	var orderSchema schemas.OrderSchema
	log := c.Locals("logger").(*slog.Logger)

	if err := c.BodyParser(&orderSchema); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	validationErrs := middleware.NewStructValidator().Validate(orderSchema)

	if len(validationErrs) > 0 {
		return fiber.NewError(fiber.StatusBadRequest, strings.Join(validationErrs, ","))
	}
	order := models.Order{Status: string(orderSchema.Status)}
	var totalAmount float64
	for _, item := range orderSchema.OrderItems {
		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
		totalAmount = float64(item.Quantity) * item.UnitPrice
		order.OrderItems = append(order.OrderItems, orderItem)
	}
	order.TotalAmount = totalAmount

	// Save the order to the database
	db := c.Locals("db").(*gorm.DB)
	result := db.Save(&order)
	if result.Error != nil {
		log.Error("Failed to create order in the database", "error", result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create order"+result.Error.Error())
	}
	go middleware.PublishOrder(&order, log)
	return c.Status(fiber.StatusOK).JSON(order)
}

// Update Order
//
//	@Summary		Updaate Order
//	@Description	Update an existing order by ID
//	@Tags			Orders
//	@Produce		json
//	@Accept			json
//
//	@param			id		path		int							true	"Order ID"
//
//	@Param			order	body		schemas.OrderUpdateSchema	true	"Order payload"
//	@Success		200		{object}	models.Order
//
//	@Failure		400		{object}	middleware.GlobalErrorHandlerResp
//	@Failure		404		{object}	middleware.GlobalErrorHandlerResp
//	@Failure		500		{object}	middleware.GlobalErrorHandlerResp
//
//	@Router			/orders/{id} [put]
func UpdateOrder(c *fiber.Ctx) error {

	log := c.Locals("logger").(*slog.Logger)
	orderId, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
	}

	var orderSchema schemas.OrderUpdateSchema
	if err := c.BodyParser(&orderSchema); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if orderSchema.Status == "" {
		log.Error("Invalid status value: empty")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid status value: empty")
	}

	log.Info("Order ID to be updated: ", "orderId", orderId)
	// Parse the request body to get updated order details

	db := c.Locals("db").(*gorm.DB)
	var order models.Order
	db.Preload("OrderItems").First(&order, orderId)
	if order.ID == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Order not found")
	}

	order.Status = string(orderSchema.Status)
	log.Info("Updated order details: ", "order", order)
	// Save the order to the database
	db.Save(&order)
	return c.Status(fiber.StatusOK).JSON(order)
}

// Fetch Order
//
//	@Summary		Fetch Order
//	@Description	Fetch an order by ID
//	@Tags			Orders
//	@Produce		json
//	@Accept			json
//
//	@param			id	path		int	true	"Order ID"
//
//	@Success		200	{object}	models.Order
//
//	@Failure		400	{object}	middleware.GlobalErrorHandlerResp
//	@Failure		404	{object}	middleware.GlobalErrorHandlerResp
//	@Failure		500	{object}	middleware.GlobalErrorHandlerResp
//
//	@Router			/orders/{id} [get]
func GetOrder(c *fiber.Ctx) error {
	log := c.Locals("logger").(*slog.Logger)
	orderId, err := c.ParamsInt("id")
	if err != nil {
		log.Error("Invalid order ID parameter", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
	}
	// Parse the request body to get updated order details

	db := c.Locals("db").(*gorm.DB)
	var order models.Order

	db.Preload("OrderItems").First(&order, uint(orderId))
	if order.ID == 0 {
		log.Error("Order not found", "orderId", orderId)
		return fiber.NewError(fiber.StatusNotFound, "Order not found")
	}
	return c.Status(fiber.StatusOK).JSON(OrderResponse{Order: order, User: populateUserDetails(order.UserId)})
}

// Delete Order
//
//	@Summary		Delete Order
//	@Description	Delete an order by ID
//	@Tags			Orders
//	@Produce		json
//	@Accept			json
//
//	@param			id	path	int	true	"Order ID"
//
//	@Success		204	"Order deleted successfully"
//
//	@Failure		400	{object}	middleware.GlobalErrorHandlerResp
//	@Failure		404	{object}	middleware.GlobalErrorHandlerResp
//	@Failure		500	{object}	middleware.GlobalErrorHandlerResp
//
//	@Router			/orders/{id} [delete]
func DeleteOrder(c *fiber.Ctx) error {
	log := c.Locals("logger").(*slog.Logger)
	orderId, err := c.ParamsInt("id")
	if err != nil {
		log.Error("Invalid order ID parameter", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
	}
	// Parse the request body to get updated order details

	db := c.Locals("db").(*gorm.DB)
	var order models.Order
	db.First(&order, orderId)
	if order.ID == 0 {
		log.Error("Order not found", "orderId", orderId)
		return fiber.NewError(fiber.StatusNotFound, "Order not found")
	}
	db.Delete(&order)
	return c.SendStatus(fiber.StatusNoContent)
}

// Start Kafka Consumer
//
//	@Summary		Start Kafka Consumer
//	@Description	Start the Kafka consumer to process orders
//	@Tags			Orders
//	@Success		200	{object}	fiber.Map	"Kafka consumer started successfully"
//
//	@Router			/orders/consumer/start [get]
func StartConsumer(c *fiber.Ctx) error {
	middleware.StartKafkaConsumer()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Kafka consumer started successfully",
	})
}

// Stop Kafka Consumer
//
//	@Summary		Stop Kafka Consumer
//	@Description	Stop the Kafka consumer to process orders
//	@Tags			Orders
//	@Success		200	{object}	fiber.Map	"Kafka consumer stopped successfully"
//
//	@Router			/orders/consumer/stop [get]
func StopConsumer(c *fiber.Ctx) error {
	middleware.StopKafkaConsumer()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Kafka consumer stopped successfully",
	})
}

func populateUserDetails(userId uint) interface{} {
	return middleware.NewUserService().GetUser(userId)
}

type OrderResponse struct {
	Order models.Order `json:"order"`
	User  any          `json:"user"`
}
