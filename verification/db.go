package verification

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongodb struct {
	conn       *mongo.Client
	db         *mongo.Database
	Collection *mongo.Collection
	ctx        context.Context
}

// ConnDB Connect to database
// var ConnDB mongodb

func (db *mongodb) setDB(name string) {
	(*db).db = db.conn.Database(name)
}

func (db *mongodb) setCollection(name string) {
	(*db).Collection = db.db.Collection(name)
}

func (db *mongodb) setCtx() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := db.conn.Connect(ctx)
	if err != nil {
		log.Println(err)
	}
	(*db).ctx = ctx

	err = db.conn.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println(err)
	}
}

func (db *mongodb) ConnectToMongo(name string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(name))
	if err != nil {
		log.Println(err)
	}

	(*db).conn = client
}

func (db *mongodb) initConnDB(server, database, collection string) {
	db.ConnectToMongo(server)
	db.setCtx()
	db.setDB(database)
	db.setCollection(collection)
}

func (db *mongodb) close() {
	err := db.db.Client().Disconnect(db.ctx)
	if err != nil {
		log.Println(err)
	}

	err = db.conn.Disconnect(db.ctx)
	if err != nil {
		log.Println(err)
	}

	db.conn = nil
	db.ctx = nil
	db.db = nil
	db.Collection = nil
}
