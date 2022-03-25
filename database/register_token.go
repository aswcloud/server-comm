package database

import (
	"context"
	"crypto/sha512"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type RegisterTokenCollection struct {
	collection *mongo.Collection
}

func (self *Client) RegisterTokenCollection() *RegisterTokenCollection {
	user := RegisterTokenCollection{}
	if self.database == nil {
		return nil
	}

	user.collection = self.database.Collection("registerToken")
	return &user
}

func (self *RegisterTokenCollection) TokenCount() int64 {
	count, _ := self.collection.CountDocuments(context.TODO(), bson.D{})
	return count
}

// duration, Unit : Seconds
func (self *RegisterTokenCollection) CreateToken(duration int32) [64]byte {
	iat := time.Now().UTC()
	exp := iat.Add(time.Second * time.Duration(duration))

	uuid, _ := uuid.New()
	data := append(uuid[:], []byte(iat.String())...)
	data = append(data, []byte(exp.String())...)

	token := sha512.Sum512(data)

	self.collection.InsertOne(context.TODO(), bson.D{
		{"iat", iat},
		{"exp", exp},
		{"dur", duration},
		{"token", token},
	})

	return token
}

// 성공 유무와, iat, exp, dur 관련 실패 데이터
func (self *RegisterTokenCollection) ExistsToken(token []byte) (bool, error) {
	nowTime := time.Now().UTC()
	var elem bson.M

	err := self.collection.FindOne(context.TODO(), bson.D{
		{"token", token},
	}).Decode(&elem)

	if err == mongo.ErrNoDocuments {
		return false, err
	}

	exp := elem["exp"].(time.Time)
	// exp.Sub(nowTime)
	if exp.Sub(nowTime) < 0 {
		return false, fmt.Errorf("token is expired")
	}

	return true, nil
}
