package queue

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/jyotishp/kafka-test/pkg/db"
	pb "github.com/jyotishp/kafka-test/pkg/proto"
	"google.golang.org/protobuf/proto"
	"log"
	"regexp"
	"sync"
)

func NewTweetsConsumer(brokers, topics []string, group string) *TweetsConsumer {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion("2.2.1")
	if err != nil {
		log.Fatalf("unable to parse kafka version")
	}
	config.Version = version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := &TweetsConsumer{
		Ready:        make(chan bool),
		Config:       config,
		brokers:      brokers,
		topics:       topics,
		keywordsExpr: "swiggy(cares|_in)?",
	}

	consumer.ctx, consumer.ctxCancel = context.WithCancel(context.Background())
	consumer.client, err = sarama.NewConsumerGroup(brokers, group, consumer.Config)
	if err != nil {
		log.Panicf("failed to create consumer group client: %v", err)
	}

	return consumer
}

type TweetsConsumer struct {
	Ready        chan bool
	Config       *sarama.Config
	brokers      []string
	topics       []string
	client       sarama.ConsumerGroup
	ctx          context.Context
	ctxCancel    context.CancelFunc
	keywordsExpr string
}

func (c *TweetsConsumer) Consume(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if err := c.client.Consume(c.ctx, c.topics, c); err != nil {
			log.Panicf("error from consumer: %v", err)
		}

		if c.ctx.Err() != nil {
			return
		}
		c.Ready = make(chan bool)
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *TweetsConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *TweetsConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *TweetsConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for message := range claim.Messages() {
		tweet := &pb.Tweet{}
		proto.Unmarshal(message.Value, tweet)
		match, err := regexp.MatchString(c.keywordsExpr, tweet.Message)
		if err != nil {
			log.Printf("failed to parse regex: %v", err)
		}
		if match {
			db.AddTweet(db.NewDbSession(), tweet)
			log.Printf("Got an interesting tweet: %v", tweet.Message)
		}
		session.MarkMessage(message, "")
	}

	return nil
}
