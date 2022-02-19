package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

func New() *Client {
	return &Client{}
}
