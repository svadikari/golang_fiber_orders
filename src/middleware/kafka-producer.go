package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/svadikari/golang_fiber_orders/src/orders/models"
)

func PublishOrder(order *models.Order, log *slog.Logger) {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	if broker == "" {
		broker = "localhost:9092"
	}
	if topic == "" {
		topic = "orders"
	}
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		panic(err)
	}
	jsonData, err := json.Marshal(order)
	if err != nil {
		log.Error("Error marshalling struct:", "error", err.Error())
		return
	}

	wg := sync.WaitGroup{}

	// Delivery report handler for produced messages
	go func() {
		defer wg.Done()
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Error("Failed to deliver message:", "partition", ev.TopicPartition, "error", ev.TopicPartition.Error)
				} else {
					log.Info("Published order event", "orderId", string(ev.Key), "partition", ev.TopicPartition.Partition, "Offset", ev.TopicPartition.Offset)
				}
			}
		}
	}()
	wg.Add(1)
	if err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(order.ID) % 5},
		Key:            []byte(fmt.Sprintf("%d", order.ID)),
		Value:          jsonData,
	}, nil); err != nil {
		log.Error("Failed to produce message to Kafka", "error", err)
		return
	}
	p.Flush(15 * 1000)
	p.Close()
	wg.Wait()
}
