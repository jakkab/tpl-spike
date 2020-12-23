package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jakkab/tpl-spike/assets"
	"github.com/jakkab/tpl-spike/gcp"
	"github.com/jakkab/tpl-spike/kafka"
	"github.com/jakkab/tpl-spike/tpl"
	"google.golang.org/api/option"
	"log"
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

type compiler interface {
	Compile([]byte, map[string]interface{}) (string, error)
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

	ch := make(chan []byte)
	fmt.Println("Configuring Kafka consumer...")
	go kafka.ConfigureAndStartConsumer(ch, *brokerAddr, *topic, *group)
	fmt.Printf("Kafka consumer listening on %s, subscribed to topic %s", *brokerAddr, *topic)

	for msg := range ch {

		func() {

			s := new(assets.Source)
			if err := json.Unmarshal(msg, s); err != nil {
				fmt.Println("Invalid input data")
				return
			}

			tplBytes, err := s.DownloadTemplate()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			dataMap, err := s.DownloadDataSource()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var c compiler
			switch s.TemplateType {
			case "go":
				c = &tpl.GoBasic{}
			case "handlebars":
				c = &tpl.Handlebars{}
			default:
				fmt.Printf("\nUnknown template type: %s", s.TemplateType)
				return
			}

			content, err := c.Compile(tplBytes, dataMap)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if err := uploader.Do(ctx, content, fmt.Sprintf(outputFilenameFmt, time.Now().Format(time.RFC3339Nano))); err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
}
