package producing

import (
	"context"
	"log"

	"github.com/redsuperbat/kafka-commander/src/options"
	"github.com/segmentio/kafka-go"
)

func NewConn() (func(value []byte) error, func() error) {
	args := options.GetArgs()
	log.Println("Connecting to broker", args.KafkaBroker, "...")
	conn, err := kafka.DialLeader(context.Background(), "tcp", args.KafkaBroker, args.KafkaTopic, 0)
	log.Println("Connection successful!")
	if err != nil {
		log.Fatalln("failed to dial leader:", err.Error())
	}

	return func(value []byte) error {
		msg := kafka.Message{Value: value}
		if _, err := conn.WriteMessages(msg); err != nil {
			return err
		}
		return nil
	}, conn.Close

}
