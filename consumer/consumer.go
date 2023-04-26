package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func main() {
	// get kafka reader using environment variables.
	processingTime, _ := strconv.Atoi(os.Getenv("PROCESSING_TIME"))
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("GROUP_ID")

	reader := getKafkaReader(kafkaAddress, topic, groupID)

	defer reader.Close()

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		time.Sleep(time.Duration(processingTime) * time.Millisecond)
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
