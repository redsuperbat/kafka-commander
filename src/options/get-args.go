package options

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type KafkaCommanderArgs struct {
	PubKeyUrl   string
	KafkaTopic  string
	KafkaBroker string
	ServerPort  int
	CommandPath string
	LogLevel    zerolog.Level
}

var args *KafkaCommanderArgs

var flagSet flag.FlagSet = flag.FlagSet{}

func parseArgs() {
	KafkaBroker := flagSet.String("broker", "localhost:9092", "url to kafka broker")
	KafkaTopic := flagSet.String("topic", "", "topic to produce messages to")
	PubKeyUrl := flagSet.String("pub-key-url", "", "url to dynamically fetch jwt RSA public key")
	ServerPort := flagSet.Int("port", 8887, "port for api")
	CommandPath := flagSet.String("path", "/", "path for commands")
	LogLevelString := flagSet.String("log-level", "debug", "level of log printed to the console")

	err := flagSet.Parse(os.Args[1:])

	if err != nil {
		os.Exit(0)
	}

	level, err := zerolog.ParseLevel(*LogLevelString)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	args = &KafkaCommanderArgs{
		KafkaBroker: *KafkaBroker,
		KafkaTopic:  *KafkaTopic,
		PubKeyUrl:   *PubKeyUrl,
		ServerPort:  *ServerPort,
		CommandPath: *CommandPath,
		LogLevel:    level,
	}
	if args.KafkaTopic == "" {
		log.Fatal().Msg("No topic specified for kafka-commander. Please supply a topic eg. --topic test-topic ")
	}

}

func GetArgs() *KafkaCommanderArgs {
	if !flagSet.Parsed() {
		parseArgs()
	}
	return args
}
