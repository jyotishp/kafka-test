package db_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/jyotishp/kafka-test/pkg/db"
	"testing"
)

func NewDynamoDBClient() *mockDynamoDBClient {
	return &mockDynamoDBClient{
		tables: make(map[string]bool),
	}
}

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	tables map[string]bool
}

func (m *mockDynamoDBClient) CreateTable(input *dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	if _, ok := m.tables[*input.TableName]; ok {
		return nil, fmt.Errorf("table already exists")
	}
	m.tables[*input.TableName] = true
	return &dynamodb.CreateTableOutput{TableDescription: &dynamodb.TableDescription{TableName: input.TableName}}, nil
}

func (m *mockDynamoDBClient) DescribeTable(input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	return &dynamodb.DescribeTableOutput{Table: &dynamodb.TableDescription{
		TableName: input.TableName,
	}}, nil
}

func TestInitDb(t *testing.T) {
	svc := NewDynamoDBClient()
	db.InitDb(svc)
}
