package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Ping(client *mongo.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary()) // Primary DB에 대한 연결 체크

	if err != nil {
		return false
	}
	return true
}

func (self *Client) Connect() bool {
	host := os.Getenv("MONGODB_SERVER")
	port := os.Getenv("MONGODB_PORT")
	database := os.Getenv("MONGODB_DATABASE")
	username := os.Getenv("MONGODB_ID")
	password := os.Getenv("MONGODB_PW")
	log.Println("host : " + host)
	log.Println("port : " + port)
	log.Println("database : " + database)
	log.Println("username : " + username)
	log.Println("password : " + password)

	clientOptions := options.Client().ApplyURI("mongodb://" + host + ":" + port)
	if strings.ToLower(host) != "localhost" && host != "127.0.0.1" {
		clientOptions.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return false
	}

	if Ping(client) {
		fmt.Println("Database Connect Success")
	} else {
		fmt.Println("Database Connect Fail")
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
