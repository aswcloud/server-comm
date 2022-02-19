package database

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type UserCollection struct {
	collection *mongo.Collection
}

func (self *Client) GetUserCollection() *UserCollection {
	user := UserCollection{}
	if self.database == nil {
		return nil
	}

	user.collection = self.database.Collection("users")
	return &user
}

func (self UserCollection) ExistsId(userId string) bool {
	result := self.collection.FindOne(context.TODO(), bson.D{
		{"userId", userId},
	})
	var elem bson.D
	err := result.Decode(&elem)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (self UserCollection) GetUserUuid(userId string) (string, error) {
	result := self.collection.FindOne(context.TODO(), bson.D{
		{"userId", userId},
	})
	var elem bson.M
	err := result.Decode(&elem)
	if err != nil {
		return "", err
	}
	data := elem["uuid"]
	if data == nil {
		return "", fmt.Errorf("not found userId : %s", userId)
	}
	return data.(string), nil
}

func (self UserCollection) NewUser(userId, password, nickname, email string) string {
	if self.ExistsId(userId) {
		return ""
	}

	hash := sha512.Sum512([]byte(password))
	text := hex.EncodeToString(hash[:])
	b_uuid, _ := uuid.New()
	uuid := hex.EncodeToString(b_uuid[:])
	self.collection.InsertOne(context.TODO(), bson.D{
		{"uuid", uuid},
		{"userId", userId},
		{"password", text},
		{"nickname", nickname},
		{"email", email},
	})

	return uuid
}

func (self UserCollection) DeleteUser(uuid string) bool {
	result, err := self.collection.DeleteMany(context.TODO(), bson.D{
		{"uuid", uuid},
	})

	if result.DeletedCount == 0 || err != nil {
		return false
	} else {
		return true
	}
}

func (self UserCollection) Login(userId, password string) bool {
	hash := sha512.Sum512([]byte(password))
	text := hex.EncodeToString(hash[:])

	result := self.collection.FindOne(context.TODO(), bson.D{
		{"userId", userId},
		{"password", text},
	})
	var elem bson.D
	err := result.Decode(&elem)
	if err != nil {
		return false
	} else {
		return true
	}
}
