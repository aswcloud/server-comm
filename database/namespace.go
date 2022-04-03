package database

import (
	"context"
	"log"

	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Namespace struct {
	collection *mongo.Collection
}

func (self *Client) Namespace() *Namespace {
	name := Namespace{}
	if self.database == nil {
		return nil
	}

	name.collection = self.database.Collection("namespace")
	return &name
}

func (self Namespace) ListNamespace(uuid string) ([]string, error) {
	result := self.collection.FindOne(context.TODO(), bson.D{
		{"uuid", uuid},
	})
	var elem bson.M
	err := result.Decode(&elem)
	if err != nil {
		return []string{}, err
	} else {
		data := funk.Map(elem["namespace"], func(data interface{}) string {
			return data.(string)
		}).([]string)
		return data, nil
	}
}

func (self Namespace) CreateNamespace(uuid, namespace string) (string, error) {
	data := self.collection.FindOne(context.TODO(), bson.D{
		{"uuid", uuid},
	})

	t := bson.M{}
	err := data.Decode(&t)
	// 에러가 있다 == 내용물이 없다.
	if err != nil {
		result, err := self.collection.InsertOne(context.TODO(), bson.D{
			{"uuid", uuid},
			{"namespace", bson.A{
				namespace,
			}},
		})
		log.Println(result, err)
	} else {
		if funk.Contains(t["namespace"], namespace) {
			return namespace, nil
		}

		result, err := self.collection.UpdateOne(context.TODO(), bson.D{
			{"_id", t["_id"]},
		},
			bson.M{
				"$push": bson.M{
					"namespace": namespace,
				},
			},
		)
		log.Println(result, err)
	}

	// var elem bson.M

	if err != nil {
		return namespace, nil
	} else {
		return namespace, nil
	}
}
