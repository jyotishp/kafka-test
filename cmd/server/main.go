package main

import (
	"github.com/jyotishp/kafka-test/pkg/db"
	"github.com/jyotishp/kafka-test/pkg/http"
	"github.com/jyotishp/kafka-test/pkg/queue"
	"sync"
)

const (
	topic    = "tweets"
	group    = "analyzer1"
	grpcPort = 50051
	httpPort = 8080
)

var (
	brokers = []string{"localhost:9092"}
)

func main() {
	db.InitDb(db.NewDbSession())
	go http.StartGRPC(grpcPort, topic, brokers)
	go http.StartHTTP(httpPort, grpcPort)

	topics := []string{topic}
	consumer := queue.NewTweetsConsumer(brokers, topics, group)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go consumer.Consume(wg)

	<-consumer.Ready // Await till the consumer has been set up

	wg.Wait()
}
