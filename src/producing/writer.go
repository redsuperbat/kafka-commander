package producing

import (
	"context"

	"github.com/redsuperbat/kafka-commander/src/options"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func NewConn() (func(value []byte) error, func() error) {
	args := options.GetArgs()
	log.Info().Msgf("Connecting to broker %s ...", args.KafkaBroker)
	conn, err := kafka.DialLeader(context.Background(), "tcp", args.KafkaBroker, args.KafkaTopic, 0)
	log.Info().Msg("Connection successful!")
	if err != nil {
		log.Fatal().Msgf("failed to dial leader: %s", err.Error())
	}

	return func(value []byte) error {
		msg := kafka.Message{Value: value}
		if _, err := conn.WriteMessages(msg); err != nil {
			return err
		}
		return nil
	}, conn.Close

}
