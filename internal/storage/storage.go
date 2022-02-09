// package storage divided into 3 parts:
// 1. Connection utilities: connections to database
// 2. Schemas: logically extracted access to the part of specific database.
// 3. StorageProcessor: connect above 2 point into 1 user-friendly interface.
package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection is a wrapper over the mongo client, due a handle representing a
// pool of connections to a MongoDB deployment,so it's ConnectionPool
type ConnectionPool struct {
	dbName string
	client *mongo.Client
}

func NewConnPool(dbUrl string, maxPoolSize int, dbName string) ConnectionPool {
	clientOpts := options.Client().ApplyURI(dbUrl).SetMaxPoolSize(uint64(maxPoolSize))
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal("can't connect to the database: ", err)
	}
	return ConnectionPool{client: client, dbName: dbName}
}

// AccessStorage create a storage processor over the connection pool
func (cp ConnectionPool) AccessStorage() StorageProcessor {
	return StorageProcessor{
		conn:         *cp.client.Database(cp.dbName),
		inTransaction: false,
	}
}

// the storage processor is the main interaction point.
// it hold the connection to database and provide method to obtain
// different storage schema.
type StorageProcessor struct {
	conn          mongo.Database
	inTransaction bool
}

func NewStorageProcessor(conn ConnectionPool) StorageProcessor {
	return StorageProcessor{
		conn:          conn,
		inTransaction: false,
	}
}

func (sp StorageProcessor) GetMempoolSchema() MempoolSchema {
	return interface{}
}
