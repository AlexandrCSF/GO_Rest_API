package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wb_cource/internal/app/model"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewProducer(bootstrapServers, topic string) (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"client.id":         "wb-producer",
		"acks":              "all",
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

func NewConsumer(bootstrapServers, topic, groupID string) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (p *Producer) SendOrder(order *model.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: orderJSON,
		Key:   []byte(order.OrderUID),
	}

	deliveryChan := make(chan kafka.Event, 1)
	err = p.producer.Produce(message, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}

	log.Printf("Order %s sent to Kafka successfully", order.OrderUID)
	return nil
}

func (c *Consumer) ConsumeOrders(orderHandler func(*model.Order) error) error {
	for {
		msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				continue
			}
			return fmt.Errorf("failed to read message: %w", err)
		}

		var order model.Order
		err = json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("Failed to unmarshal order: %v", err)
			continue
		}

		err = orderHandler(&order)
		if err != nil {
			log.Printf("Failed to handle order %s: %v", order.OrderUID, err)
			continue
		}

		log.Printf("Order %s processed successfully", order.OrderUID)
	}
}

func (p *Producer) Close() {
	p.producer.Close()
}

func (c *Consumer) Close() {
	c.consumer.Close()
}
