// Package storage divided into 3 parts:
// 1. Connection utilities: connections to database
// 2. Schemas: logically extracted access to the part of specific database.
// 3. Conn: connect above 2 point into 1 user-friendly interface.
package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vivijj/zigo/storage/schema"
)

// ConnectionPool is a wrapper over mongo client -- a handle representing a pool of connections
//  to a MongoDB deployment and safe for concurrent use by multiple goroutines.
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
func (cp ConnectionPool) AccessStorage() Processor {
	return ProcessorFromPool(cp)
}

// Processor is the main interaction point which holds down the connection to database and provide
// method to obtain different storage schema.
type Processor struct {
	conn          mongo.Database
	session       mongo.Session
	inTransaction bool
}

// ProcessorFromPool will build the storage processor from the `connectionPool`.
func ProcessorFromPool(conn ConnectionPool) Processor {
	return Processor{
		conn:          *conn.client.Database(conn.dbName),
		inTransaction: false,
	}
}

func (p Processor) AccessCollection(collName string) *mongo.Collection {
	return p.conn.Collection(collName)
}

func (p Processor) MempoolSchema() schema.MempoolSchema {
	return schema.MempoolSchema{}
}

func (p Processor) StateSchema() schema.StateSchema {
	return schema.StateSchema{}
}

func (p Processor) BlockSchema() schema.BlockSchema {
	return schema.BlockSchema{}
}

func (p Processor) OperationSchema() schema.OperationSchema {
	return schema.OperationSchema{}
}
