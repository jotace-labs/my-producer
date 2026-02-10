package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/orasis-holding/pricing-go-swiss-army-lib/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
    log.Print("starting producer")
	ctx := context.Background()

	metrics := kafka.NewCommonMetrics(prometheus.NewRegistry(), "trace")

    // Get Kafka brokers from environment variable, default to redpanda service
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		// kafkaBrokers = "redpanda.redpanda.svc.cluster.local:9092"
        log.Fatal("kafka brokers is empty")
	}

	fmt.Println("address: ", kafkaBrokers)

    // Create the producer
    producer, err := kafka.NewProducer(
        ctx,
        []kafka.ClientOption{
            kafka.WithKafkaClientOpts(
                kgo.SeedBrokers(kafkaBrokers), // // only works with port forwarding
                // kgo.SeedBrokers("192.168.49.2:30092"),
				
                // Other franz-go options...
            ),
        },
        kafka.WithProducerName("my-producer"),
        kafka.WithProducerMetrics(metrics), // Optional
    )
    if err != nil {
        log.Fatal("Failed to create producer:", err)
    }

    log.Print("connected. producing")

    counter := 0

    for {
        ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*30)
        defer cancel()

        message := "Hello World " + strconv.Itoa(counter)

        counter += 1

        err = producer.ProduceMessages(ctxTimeout, 
            &kgo.Record{
                Topic: "test-topic",
                Key:   []byte("key1"),
                Value: []byte(message),
            },
        )
        if err != nil {
            log.Fatal("failed to send message:", err)
            
	    } else {
            log.Print("message produced")
        }

        time.Sleep(30*time.Second)
    }

    

   
}