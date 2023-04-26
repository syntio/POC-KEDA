package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	// Load env variables
	processingTime := 1200              // strconv.Atoi(os.Getenv("PROCESSING_TIME"))
	kafkaAddress := "20.31.76.141:9093" // os.Getenv("KAFKA_ADDRESS")
	topic := "keda-topic-p3"            // os.Getenv("KAFKA_TOPIC")

	// Configure the producer
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// Create a new sync producer
	client, err := sarama.NewClient([]string{kafkaAddress}, config)
	if err != nil {
		log.Fatalln("Failed to create a client:", err)
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Fatalln("Failed to create producer:", err)
	}
	defer producer.Close()

	i := 0
	ticker := time.NewTicker(time.Duration(processingTime) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			i++

			// Create a new message to send
			message := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder("message ID: " + string(rune(i))),
			}

			// Send the message
			partition, offset, err := producer.SendMessage(message)
			if err != nil {
				log.Fatalln("Failed to send message:", err)
			}

			// Print the message content
			fmt.Printf("Message with id=%d, sent to partition %d at offset %d\n", i, partition, offset)
		}
	}
}
