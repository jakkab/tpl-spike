package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func ConfigureAndStartConsumer(c chan []byte, brokerAddr, topic, groupID string) {

	cfg := kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		Topic:    topic,
		GroupID:  groupID,
		MaxBytes: 100,
	}

	reader := kafka.NewReader(cfg)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("failed to read a message: %s", err.Error())
			continue
		}

		c <- msg.Value
	}
}
