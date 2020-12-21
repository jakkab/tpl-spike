package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jakkab/tpl-spike/kafka"
	"github.com/jakkab/tpl-spike/template"
	"io"
	"log"
	"os"
	"time"
)

const outputFilenameFmt = "compiled-%s.html"

var (
	brokerAddr = flag.String("kafka", "", "kafka broker address")
	topic      = flag.String("topic", "", "kafka topic")
	group      = flag.String("group", "my-group-1", "kafka group")
	bucket     = "tpl-spike"
)

func main() {

	flag.Parse()

	if *brokerAddr == "" {
		log.Fatal("kafka broker address not provided (e.g. --brokerAddr=localhost:9092)")
	}

	if *topic == "" {
		log.Fatal("kafka topic not provided (e.g. --topic=my-test-topic)")
	}

	//storageClient, err := storage.NewClient(context.Background(), option.WithCredentialsFile())
	//if err != nil {
	//	log.Fatal("Unable to init GCP storage client")
	//}

	fmt.Println("Configuring Kafka consumer...")

	c := make(chan []byte)
	go kafka.ConfigureAndStartConsumer(c, *brokerAddr, *topic, *group)

	fmt.Printf("Kafka consumer listening on %s, subscribed to topic %s", *brokerAddr, *topic)

	for msg := range c {

		s := new(template.Source)
		if err := json.Unmarshal(msg, s); err != nil {
			fmt.Println("Invalid input data")
			continue
		}

		fmt.Printf("\nJson comes from %s", s.JSONDataURL)
		fmt.Printf("\nTemplate comes from %s", s.TemplateURL)

		outputFilename := fmt.Sprintf(outputFilenameFmt, time.Now().String())
		outputFile, err := os.Create(outputFilename)
		if err != nil {
			log.Fatalf("unable to create file: %s, %s", outputFilename, err.Error())
		}

		if err := s.Compile(outputFile); err != nil {
			log.Fatal(err)
		}

		file, err := os.Open(outputFilename)
		if err != nil {
			log.Fatal(err)
		}

		io.Copy(os.Stdout, file)
	}
}
