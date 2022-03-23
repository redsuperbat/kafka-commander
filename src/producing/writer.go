package producing

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

const (
	KAFKA_BROKER = "KAFKA_BROKER"
	KAFKA_TOPIC  = "KAFKA_TOPIC"
)

func NewConn() (func(value []byte) error, func() error) {
	kafkaBroker := os.Getenv(KAFKA_BROKER)
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	kafkaTopic := os.Getenv(KAFKA_TOPIC)
	if kafkaTopic == "" {
		log.Fatal().Msg("A topic must be provided for kafka commander")
	}

	log.Info().Msgf("Connecting to broker %s ...", kafkaBroker)
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBroker, kafkaTopic, 0)

	if err != nil {
		log.Fatal().Msgf("failed to dial leader: %s", err.Error())
	}

	log.Info().Msg("Connection successful!")

	return func(value []byte) error {
		msg := kafka.Message{Value: value}
		if _, err := conn.WriteMessages(msg); err != nil {
			return err
		}
		return nil
	}, conn.Close

}
