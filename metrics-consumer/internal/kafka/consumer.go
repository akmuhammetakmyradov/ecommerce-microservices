package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	sessionTimeoutMs   = 300000
	heartbeatInterval  = 100000
	maxPollInterval    = 600000
	autoCommitInterval = 5000
)

type Handler interface {
	HandleMessage(message []byte, topic kafka.Offset) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
}

func NewConsumer(handler Handler, address []string, topic, consumerGroup string) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeoutMs,
		"heartbeat.interval.ms":    heartbeatInterval,
		"max.poll.interval.ms":     maxPollInterval,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  autoCommitInterval,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	if err := c.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		handler:  handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer context cancelled, stopping...")
			return nil
		default:
			ev := c.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch msg := ev.(type) {
			case *kafka.Message:
				if err := c.handler.HandleMessage(msg.Value, msg.TopicPartition.Offset); err != nil {
					log.Printf("Error handling message: %v", err)
					continue
				}

				if _, err := c.consumer.StoreMessage(msg); err != nil {
					log.Printf("Error storing message offset: %v", err)
					continue
				}

			case kafka.Error:
				log.Printf("Kafka error: %v", msg)
				if msg.Code() == kafka.ErrAllBrokersDown {
					log.Println("All brokers are down")
				}
			default:
				log.Printf("Unhandled event type: %T", msg)
			}
		}
	}
}

func (c *Consumer) Stop() error {
	log.Println("Committing and closing consumer...")
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	return c.consumer.Close()
}
