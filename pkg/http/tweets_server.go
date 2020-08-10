package http

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	pb "github.com/jyotishp/kafka-test/pkg/proto"
	"github.com/jyotishp/kafka-test/pkg/queue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Tweets server
type TweetsServer struct {
	pb.UnimplementedTweetsServiceServer
	kafkaTopic string
	collector  sarama.SyncProducer
}

// Create a new instance of tweets server
func NewTweetsServer(kafkaTopic string, brokers []string) *TweetsServer {
	return &TweetsServer{
		kafkaTopic: kafkaTopic,
		collector:  queue.NewKafkaProducer(brokers),
	}
}

// Create a new tweet and push it to kafka
func (s *TweetsServer) AddTweet(ctx context.Context, i *pb.Tweet) (*pb.Tweet, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create uuid: %v", err)
	}
	tweet := &pb.Tweet{
		TweetId: id.String(),
		Message: i.Message,
	}
	err = s.PushToKafka(tweet)
	if err != nil {
		return nil, err
	}
	return tweet, nil
}

// Push the tweet to kafka stream
func (s *TweetsServer) PushToKafka(tweet *pb.Tweet) error {
	msg, err := proto.Marshal(tweet)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to marshal tweet: %v", err)
	}
	_, _, err = s.collector.SendMessage(&sarama.ProducerMessage{
		Topic: s.kafkaTopic,
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to inject to kafka: %v", err)
	}
	return nil
}
