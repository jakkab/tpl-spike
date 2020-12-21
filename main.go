package main

import (
	"flag"
	"fmt"
	"github.com/jakkab/tpl-spike/kafka"
	"log"
)

var (
	brokerAddr = flag.String("kafka", "", "kafka broker address")
	topic      = flag.String("topic", "", "kafka topic")
	group      = flag.String("group", "my-group-1", "kafka group")
)

func main() {

	flag.Parse()

	if *brokerAddr == "" {
		log.Fatal("kafka broker address not provided (e.g. --brokerAddr=localhost:9092)")
	}

	if *topic == "" {
		log.Fatal("kafka topic not provided (e.g. --topic=my-test-topic)")
	}

	fmt.Printf("Kafka consumer listening on %s:", *brokerAddr)

	c := make(chan string)
	go kafka.ConfigureAndStartConsumer(c, *brokerAddr, *topic, *group)

	fmt.Println("Kafka consumer up and running")

	for msg := range c {

		

	}
}
