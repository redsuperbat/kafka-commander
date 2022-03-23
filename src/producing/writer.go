package producing

import (
	"context"
	"math"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

const (
	KAFKA_BROKER = "KAFKA_BROKER"
	KAFKA_TOPIC  = "KAFKA_TOPIC"

	ceilRetryTime float64 = 25
)

func getRetryTime(numberOfRetries int64) float64 {
	if numberOfRetries > 5 {
		return 25
	}
	retryTime := math.Exp(float64(numberOfRetries))
	return retryTime * 100
}

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
	var numberOfRetries int64 = 0
	for err == nil {

		conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBroker, kafkaTopic, 0)

		if err != nil {
			numberOfRetries += 1
			retryTime := getRetryTime(numberOfRetries)
			log.Error().Msgf("failed to dial leader: %s. Retrying in %v seconds", err.Error(), retryTime/1000)
			time.Sleep(time.Duration(retryTime) * time.Millisecond)
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
