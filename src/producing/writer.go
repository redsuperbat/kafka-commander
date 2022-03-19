package producing

import (
	"time"

	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(kafkaURL),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
}
