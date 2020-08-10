package http

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama/mocks"
	pb "github.com/jyotishp/kafka-test/pkg/proto"
	"github.com/jyotishp/kafka-test/pkg/queue"
	"testing"
)

func TestTweetsServer_AddTweet(t *testing.T) {
	producer := mocks.NewSyncProducer(t, queue.ProducerConfig())
	producer.ExpectSendMessageAndSucceed()
	s := &TweetsServer{
		kafkaTopic: "topic",
		collector:  producer,
	}
	ctx := context.Background()
	input := &pb.Tweet{Message: "Some random msg"}
	_, err := s.AddTweet(ctx, input)
	if err != nil {
		t.Errorf("failed to add tweet: %v", err)
	}

	producer.ExpectSendMessageAndFail(fmt.Errorf("failed to inject to kafka"))
	s.collector = producer
	_, err = s.AddTweet(ctx, input)
	if err == nil {
		t.Errorf("expected error, but got no error")
	}
}
