package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Host      string
	Port      int
	connected bool
}

var Mongo *mongo.Client

func (db *Database) Connect() {
	if db.connected == true {
		log.Println("already connected to database.")
		return
	}
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", db.Host, db.Port))

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	db.connected = true
	Mongo = client

	fmt.Println("Connected to MongoDB")
}

func (db *Database) Save(database string, collection string, data ...interface{}) {
	c := Mongo.Database(database).Collection(collection)
	insertR, err := c.InsertMany(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted documents: ", insertR.InsertedIDs)
}

func (db *Database) FindOne(database string, collection string, filter bson.D) interface{} {
	var result interface{}

	c := Mongo.Database(database).Collection(collection)
	err := c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
