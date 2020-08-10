package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	pb "github.com/jyotishp/kafka-test/pkg/proto"
)

// Create a new tweet entity and add it to the database
func AddTweet(svc dynamodbiface.DynamoDBAPI, tweet *pb.Tweet) error {
	err := AddObject(svc, TweetsTable, tweet)
	if err != nil {
		return err
	}
	return nil
}

// Add an object to the database
func AddObject(svc dynamodbiface.DynamoDBAPI, table string, data interface{}) error {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return fmt.Errorf("failed to marshal new customer: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to add item to database: %v", err)
	}
	return nil
}
