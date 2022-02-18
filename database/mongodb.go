package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient() (*mongo.Client, error) {
	// 몽고DB 연결
	clientOptions := options.Client().ApplyURI("mongodb://localhost:20000")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("몽고 DB에 연결했습니다!")

	// 내용을 적을 부분

	// 몽고DB 연결 끊기
	uesrsCollection := client.Database("servercomm").Collection("users")

	uesrsCollection.DeleteMany(context.TODO(), bson.D{{"userID", "test"}})

	for i := 0; i < 10; i++ {
		uesrsCollection.InsertOne(context.TODO(), bson.D{
			{"userID", "test"},
			{"uuid", i},
		})
	}
	cursor, _ := uesrsCollection.Find(context.TODO(), bson.D{
		{"uuid", bson.D{{"$gte", 5}}},
	})

	for cursor.Next(context.TODO()) {
		var elem bson.M
		err := cursor.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}
		// find 결과 print
		fmt.Println(elem)
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("몽고DB 연결을 종료했습니다!")

	return client, nil
}
