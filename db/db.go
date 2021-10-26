package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Db struct {
	client *mongo.Client
	ctx context.Context
	portsDb *mongo.Database
	ipsCollection *mongo.Collection
}

func (db *Db) Connect() error {
	uri := "mongodb://127.0.0.1:27017"
	var err error = nil
	db.client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	db.ctx, _ = context.WithCancel(context.Background())
	err = db.client.Connect(db.ctx)
	if err != nil {
		log.Fatal(err)
	}

	db.portsDb = db.client.Database("ports")
	db.ipsCollection = db.portsDb.Collection("ips")

	return nil
}

func (db *Db) Disconnect() {
	err := db.client.Disconnect(db.ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *Db) ListDbs() {
	databases, err := db.client.ListDatabaseNames(db.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}

func (db *Db) InsertRow(ip string, port string) {
	_, err := db.ipsCollection.InsertOne(db.ctx, bson.D{
		{Key: "ip", Value: ip},
		{Key: "port", Value: port},
	})
	if err != nil {
		log.Println("Insert failed: ", err)
	}
}