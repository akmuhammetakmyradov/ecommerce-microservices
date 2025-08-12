package kafka

import (
	"cart/internal/repository/interfaces"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeout       = 5000
	partition    int32 = 0
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer(address []string, topic string) (interfaces.KafkaProd, error) {
	if len(address) == 0 {
		return nil, fmt.Errorf("kafka broker address list is empty")
	}

	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}

	prod, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka producer: %w", err)
	}

	go func() {
		for e := range prod.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("❌ Delivery failed: %v", ev.TopicPartition)
				} else {
					log.Printf("✅ Delivered to %v", ev.TopicPartition)
				}
			case kafka.Error:
				log.Printf("Kafka error: %v", ev)
			}
		}
	}()

	return &Producer{producer: prod, topic: topic}, nil
}

func (p *Producer) Produce(message []byte, key string, t time.Time) error {
	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: partition,
		},
		Value:     message,
		Key:       []byte(key),
		Timestamp: t,
	}

	if err := p.producer.Produce(kafkaMessage, nil); err != nil {
		return fmt.Errorf("error sending message to kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
