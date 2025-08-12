package interfaces

import "time"

type KafkaProd interface {
	Produce(message []byte, key string, t time.Time) error
	Close()
}
