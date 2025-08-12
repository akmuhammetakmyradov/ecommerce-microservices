package service

import (
	"encoding/json"
	"stocks/internal/models"
	"time"
)

type Payload struct {
	SKU   uint32 `json:"sku"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

type KafkaEvent struct {
	Type      string    `json:"type"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
	Payload   Payload   `json:"payload"`
}

func BuildKafkaEvent(addedType string, item models.StockItem) ([]byte, time.Time, error) {
	timestamp := time.Now()
	message := KafkaEvent{
		Type:      addedType,
		Service:   "stock",
		Timestamp: timestamp,
		Payload: Payload{
			SKU:   item.SKU,
			Count: item.Count,
			Price: item.Price,
		},
	}

	msg, err := json.Marshal(message)
	if err != nil {
		return nil, time.Time{}, err
	}

	return msg, timestamp, nil
}
