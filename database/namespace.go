package database

import (
	"context"
	"log"

	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Namespace struct {
	collection      *mongo.Collection
	totalCollection *mongo.Collection
}

func (self *Client) Namespace() *Namespace {
	name := Namespace{}
	if self.database == nil {
		return nil
	}
	name.totalCollection = self.database.Collection("totalNamespace")
	name.collection = self.database.Collection("namespace")
	return &name
}

func (self Namespace) Exists(namespace string) bool {
	result := self.totalCollection.FindOne(context.TODO(), bson.D{
		{"namespace", namespace},
	})
	var elem bson.M
	err := result.Decode(&elem)
	if err != nil {
		return false
	} else {
		return true
	}
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

func (self Namespace) CreateNamespace(uuid, namespace string) (string, bool) {
	data := self.collection.FindOne(context.TODO(), bson.D{
		{"uuid", uuid},
	})
	self.totalCollection.InsertOne(context.TODO(), bson.D{
		{"namespace", namespace},
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
			return namespace, true
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
	return namespace, false
}

func (self Namespace) DeleteNamespace(uuid, namespace string) (string, bool) {
	self.totalCollection.DeleteMany(context.TODO(), bson.D{
		{"namespace", namespace},
	})

	result, err := self.collection.UpdateOne(context.TODO(), bson.D{
		{"uuid", uuid},
	},
		bson.M{
			"$pull": bson.M{
				"namespace": namespace,
			},
		},
	)
	log.Println(result, err)
	return namespace, true
}
