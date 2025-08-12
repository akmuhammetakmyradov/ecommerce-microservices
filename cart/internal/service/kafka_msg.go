package service

import (
	"cart/internal/models"
	"encoding/json"
	"time"
)

type Payload struct {
	CardId int64  `json:"card_id"`
	SKU    uint32 `json:"sku"`
	Count  uint32 `json:"count"`
	Price  uint32 `json:"price"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type KafkaEvent struct {
	Type      string    `json:"type"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
	Payload   Payload   `json:"payload"`
}

func BuildKafkaEvent(addedType string, cartId int64, price uint32, reason, status string, item models.CartItem) ([]byte, time.Time, error) {
	timestamp := time.Now()
	message := KafkaEvent{
		Type:      addedType,
		Service:   "cart",
		Timestamp: timestamp,
		Payload: Payload{
			CardId: cartId,
			SKU:    item.SKU,
			Count:  item.Count,
			Price:  price,
			Reason: reason,
			Status: status,
		},
	}

	msg, err := json.Marshal(message)
	if err != nil {
		return nil, time.Time{}, err
	}

	return msg, timestamp, nil
}
