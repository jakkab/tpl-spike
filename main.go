package main

import (
	"fmt"
	consumer "github.com/jakkab/tpl-spike/kafka"
)

var (
	brokerAddr = "kafka.default.svc.cluster.local:9092"
	topic = "my-topic-1"
	group = "my-group-1"
)

func main() {

	fmt.Println("Kafka consumer starting...")

	c := make(chan string)
	go consumer.ConfigureAndStartConsumer(c, brokerAddr, topic, group)

	fmt.Println("Kafka consumer up and running")

	for msg := range c {
		fmt.Println(msg)
	}

	/// Part 1, main microservice

}
