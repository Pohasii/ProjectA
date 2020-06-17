package verification

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"strconv"
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

	db.conn = nil
	db.ctx = nil
	db.db = nil
	db.Collection = nil
}

func checkAkk(connDb mongodb, Credential Credentials) bool {

	filter := bson.D{
		{"password", Credential.Password},
		{"login", Credential.Login},
	}

	count, err := connDb.Collection.CountDocuments(connDb.ctx, filter)
	if err != nil {
		log.Println("auth service, checkAkk error: ", err)
	}

	if count == 1 {
		return true
	}
	return false
}

func updateAkk(connDb mongodb, Credential Credentials) (Client Profile, status bool) {

	h := md5.New()
	time := strconv.Itoa(int(time.Now().Unix()))
	io.WriteString(h, time)
	token := hex.EncodeToString(h.Sum([]byte(time)))

	update := bson.D{
		{"$set", bson.D{
			{"token", token},
		}},
	}

	filter := bson.D{
		{"password", Credential.Password},
		{"login", Credential.Login},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true)

	err := connDb.Collection.FindOneAndUpdate(connDb.ctx, filter, update, opts).Decode(&Client)
	if err != nil {
		log.Println("auth service, updateAkk error: ", err)
		return Client, false
	}

	return Client, true
}
