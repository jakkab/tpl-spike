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
)

var (
	brokerAddr      = flag.String("kafka", "", "kafka broker address")
	topic           = flag.String("topic", "", "kafka topic")
	group           = flag.String("group", "my-group-1", "kafka group")
	bucket          = "tpl-spike"
	credentialsFile = "/etc/sa/sa_key.json"
)

type Source struct {
	TemplateURL string `json:"templateURL"`
	JSONDataURL string `json:"jsonDataURL"`
}

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
	uploader := gcp.NewGcpUploader(storageClient, bucket)

	c := make(chan []byte)
	fmt.Println("Configuring Kafka consumer...")
	go kafka.ConfigureAndStartConsumer(c, *brokerAddr, *topic, *group)
	fmt.Printf("Kafka consumer listening on %s, subscribed to topic %s", *brokerAddr, *topic)

	for msg := range c {

		func() {
			s := new(Source)
			if err := json.Unmarshal(msg, s); err != nil {
				fmt.Println("Invalid input data")
				return
			}

			fmt.Printf("\nJson comes from %s", s.JSONDataURL)
			fmt.Printf("\nTemplate comes from %s", s.TemplateURL)

			bc := &template.GoBasic{}

			outputFilename, err := bc.Compile(s.TemplateURL, s.JSONDataURL)
			if err != nil {
				fmt.Println(err)
				return
			}

			if err := uploader.Do(ctx, outputFilename); err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
}
