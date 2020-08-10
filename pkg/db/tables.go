package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"log"
)

const TweetsTable = "Tweets"

var (
	tweetId = &dynamodb.KeySchemaElement{
		AttributeName: aws.String("tweet_id"),
		KeyType:       aws.String("HASH"),
	}
	tweetIdAttr = &dynamodb.AttributeDefinition{
		AttributeName: aws.String("tweet_id"),
		AttributeType: aws.String("S"),
	}
)

// Create a new table
func CreateTable(svc dynamodbiface.DynamoDBAPI, table *dynamodb.CreateTableInput) error {
	out, err := svc.CreateTable(table)
	if err != nil {
		return fmt.Errorf("unable to create table: %v", err)
	}
	log.Printf("created table: %v", out)
	return nil
}

// Check if the given table exists
func TableExists(svc dynamodbiface.DynamoDBAPI, table string) bool {
	log.Printf("checking if table %v exists", table)
	_, err := svc.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(table)})
	if err != nil {
		return false
	}
	return true
}

// Create a table if it doesn't exist
func CreateTableIfNotExists(svc dynamodbiface.DynamoDBAPI, table *dynamodb.CreateTableInput) error {
	if !TableExists(svc, *table.TableName) {
		log.Printf("creating table: %v", *table.TableName)
		return CreateTable(svc, table)
	}
	return nil
}

// Schema for tweets table
func TweetsTableSchema() *dynamodb.CreateTableInput {
	table := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{tweetIdAttr},
		KeySchema:            []*dynamodb.KeySchemaElement{tweetId},
		TableName:            aws.String(TweetsTable),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	return table
}
