package middleware

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/svadikari/golang_fiber_orders/src/orders/models"
)

type consumeOrders struct {
	run bool
}

var consumerInstance *consumeOrders

func init() {
	consumerInstance = &consumeOrders{run: false}

}

func StartKafkaConsumer() {
	if !consumerInstance.run {
		consumerInstance.run = true
		go consumerInstance.ConsumeOrders()
	} else {
		slog.Warn("Kafka consumer is already running")
	}
}

func StopKafkaConsumer() {
	if consumerInstance.run {
		consumerInstance.run = false
		slog.Info("Kafka consumer stopped")
	} else {
		slog.Warn("Kafka consumer is not running")
	}
}

func (con *consumeOrders) ConsumeOrders() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	consumerGroup := os.Getenv("KAFKA_CONSUMER_GROUP")

	if broker == "" {
		broker = "localhost:9092"
	}
	if topic == "" {
		topic = "orders"
	}

	if consumerGroup == "" {
		consumerGroup = "order-consumer-group"
	}

	p, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          consumerGroup,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		slog.Error("Failed to create Kafka consumer", "error", err)
		return
	}

	if err := p.SubscribeTopics([]string{topic}, nil); err != nil {
		slog.Error("Failed to subscribe to Kafka topic", "error", err)
		return
	}
	slog.Info("Kafka consumer started, listening to topic:", "topic", topic)

	for con.run {
		msg, err := p.ReadMessage(-1)
		if err == nil {
			var order models.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				slog.Error("Failed to unmarshal Kafka message", "error", err)
				continue
			}
			slog.Info("Consumed order from Kafka", "key", msg.Key, "order", order)
			// Process the order as needed, e.g., update database, trigger other actions, etc.
		} else {
			// The client will automatically try to recover from all errors.
			slog.Error("Consumer error", "error", err)
		}
	}
}
