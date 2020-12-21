package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jakkab/tpl-spike/gcp"
	"github.com/jakkab/tpl-spike/kafka"
	"github.com/jakkab/tpl-spike/template"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

const outputFilenameFmt = "compiled-%s.html"

var (
	brokerAddr      = flag.String("kafka", "", "kafka broker address")
	topic           = flag.String("topic", "", "kafka topic")
	group           = flag.String("group", "my-group-1", "kafka group")
	bucket          = "tpl-spike"
	credentialsFile = "/etc/sa/sa_key.json"
)

func main() {

	flag.Parse()

	if *brokerAddr == "" {
		log.Fatal("kafka broker address not provided (e.g. --brokerAddr=localhost:9092)")
	}

	if *topic == "" {
		log.Fatal("kafka topic not provided (e.g. --topic=my-test-topic)")
	}

	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatal("Unable to init GCP storage client")
	}

	uploader := gcp.NewUploader(storageClient)

	fmt.Println("Configuring Kafka consumer...")

	c := make(chan []byte)
	go kafka.ConfigureAndStartConsumer(c, *brokerAddr, *topic, *group)

	fmt.Printf("Kafka consumer listening on %s, subscribed to topic %s", *brokerAddr, *topic)

	for msg := range c {

		func() {
			s := new(template.Source)
			if err := json.Unmarshal(msg, s); err != nil {
				fmt.Println("Invalid input data")
				return
			}

			fmt.Printf("\nJson comes from %s", s.JSONDataURL)
			fmt.Printf("\nTemplate comes from %s", s.TemplateURL)

			outputFilename := fmt.Sprintf(outputFilenameFmt, time.Now().String())
			outputFile, err := os.Create(outputFilename)
			if err != nil {
				fmt.Printf("unable to create file: %s, %s", outputFilename, err.Error())
				return
			}

			if err := s.Compile(outputFile); err != nil {
				fmt.Println(err)
				return
			}

			if err := uploader.Do(ctx, outputFilename, bucket); err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
}
