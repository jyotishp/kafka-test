package queue

import (
	"github.com/Shopify/sarama"
	"log"
)

func NewKafkaProducer(brokers []string) sarama.SyncProducer {
	config := ProducerConfig()
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("failed to connect to kafka: %v", err)
	}
	return producer
}

func ProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	return config
}
