package producing

import (
	"context"
	"os"
	"time"

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
	var err error
	numberOfRetries := 0
	for err == nil {

		conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBroker, kafkaTopic, 0)

		if err != nil {
			log.Error().Msgf("failed to dial leader: %s. Retrying in 3 seconds", err.Error())
			time.Sleep(3 * time.Second)
			if numberOfRetries > 5 {
				log.Fatal().Msgf("Failed to dial leader 5 times, exiting")
			}
			numberOfRetries += 1
			continue
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
	log.Fatal().Msg("Unable to connect to kafka, please try again..")
	return nil, nil
}
