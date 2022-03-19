package options

import (
	"flag"
	"log"
)

type KafkaCommanderArgs struct {
	PubKeyUrl   string
	KafkaTopic  string
	KafkaBroker string
	ServerPort  int
	CommandPath string
}

var args *KafkaCommanderArgs

func parseArgs() {
	KafkaBroker := flag.String("broker", "localhost:9092", "url to kafka broker")
	KafkaTopic := flag.String("topic", "", "topic to produce messages to")
	PubKeyUrl := flag.String("pub-key-url", "", "url to dynamically fetch jwt RSA public key")
	ServerPort := flag.Int("port", 8887, "port for api")
	CommandPath := flag.String("path", "/", "path for commands")
	flag.Parse()
	args = &KafkaCommanderArgs{
		KafkaBroker: *KafkaBroker,
		KafkaTopic:  *KafkaTopic,
		PubKeyUrl:   *PubKeyUrl,
		ServerPort:  *ServerPort,
		CommandPath: *CommandPath,
	}
	if args.KafkaTopic == "" {
		log.Fatalln("No topic specified for kafka-commander. Please supply a topic eg. --topic test-topic ")
	}

}

func GetArgs() *KafkaCommanderArgs {
	log.Println("Getting args...")
	if args == nil {
		log.Println("Args not cached, parsing them...")
		parseArgs()
	}
	return args
}
