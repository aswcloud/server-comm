package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Client) Connect() bool {
	host := os.Getenv("MONGODB_SERVER")
	port := os.Getenv("MONGODB_PORT")
	database := os.Getenv("MONGODB_DATABASE")
	clientOptions := options.Client().ApplyURI("mongodb://" + host + ":" + port)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return false
	}
	self.client = client
	self.database = client.Database(database)

	return true
}

func (self *Client) Disconnect() bool {
	if self.client != nil {
		err := self.client.Disconnect(context.TODO())
		if err == nil {
			return true
		}
	}

	return false
}
