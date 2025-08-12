package main

import (
	"context"
	"log"
	"metrics-consumer/internal/config"
	"metrics-consumer/internal/handler"
	kconstructor "metrics-consumer/internal/kafka"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.GetConfig()
	h := handler.NewHandler()

	c, err := kconstructor.NewConsumer(h, cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	log.Println("✅ Successfully created metrics consumer")

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := c.Start(ctx); err != nil {
			log.Printf("Consumer exited with error: %v", err)
		}
	}()

	sig := <-sigChan
	log.Printf("🛑 Shutdown signal received: %s", sig)

	cancel()

	if err := c.Stop(); err != nil {
		log.Printf("⚠️ Error while stopping consumer: %v", err)
	} else {
		log.Println("✅ Consumer stopped gracefully")
	}
}
