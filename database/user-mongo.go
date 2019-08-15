package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	host      string
	port      int
	connected bool
	mongo     *mongo.Client
}

type Collection struct {
	name string
}

func (db *Database) Connect() {
	if db.connected == true {
		log.Println("already connected to database.")
		return
	}
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", db.host, db.port))

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	db.mongo = client
	db.connected = true

	fmt.Println("Connected to MongoDB")
}

func (db *Database) Save(database string, collection string, data ...interface{}) {
	c := db.mongo.Database(database).Collection(collection)
	insertR, err := c.InsertMany(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted documents: ", insertR.InsertedIDs)
}
