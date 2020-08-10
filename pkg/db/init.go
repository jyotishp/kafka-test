package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/jyotishp/kafka-test/pkg/utils"
	"log"
)

// Create a new AWS session
func NewDbSession() *dynamodb.DynamoDB {
	dbEndpoint := utils.GetEnv("DB_ENDPOINT", "localhost:8000")
	session, err := session.NewSession(&aws.Config{
		Endpoint:   aws.String(dbEndpoint),
		DisableSSL: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("failed to create new session: %v", err)
	}
	svc := dynamodb.New(session)
	return svc
}

// Initialize the database.
// Create new DynamoDB service and initializes the tables
func InitDb(svc dynamodbiface.DynamoDBAPI) {
	log.Printf("initializing db...")
	err := CreateTableIfNotExists(svc, TweetsTableSchema())
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}
}
