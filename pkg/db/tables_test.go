package db_test

import (
	"github.com/jyotishp/kafka-test/pkg/db"
	"testing"
)

func TestTweetsTableSchema(t *testing.T) {
	table := db.TweetsTableSchema()
	assertTableName(t, *table.TableName, db.TweetsTable)
}

func assertTableName(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("table name doesn't match")
	}
}

func TestCreateTable(t *testing.T) {
	svc := NewDynamoDBClient()
	err := db.CreateTable(svc, db.TweetsTableSchema())
	if err != nil {
		t.Errorf("unable to create new table: %v", err)
	}
}

func TestTableExists(t *testing.T) {
	svc := NewDynamoDBClient()
	db.CreateTable(svc, db.TweetsTableSchema())
	res := db.TableExists(svc, db.TweetsTable)
	if !res {
		t.Errorf("table exists but returns doesn't exist")
	}
}

func TestCreateTableIfNotExists(t *testing.T) {
	svc := NewDynamoDBClient()
	// create table
	err := db.CreateTableIfNotExists(svc, db.TweetsTableSchema())
	if err != nil {
		t.Errorf("%v", err)
	}
	// should not create table
	err = db.CreateTableIfNotExists(svc, db.TweetsTableSchema())
	if err != nil {
		t.Errorf("%v", err)
	}
}
