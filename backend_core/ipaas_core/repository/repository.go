package repository

import (
	"context"
	"encoding/json"
	"errors"

	"fermion/backend_core/db"
	"fermion/backend_core/pkg/util/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IpaasFeatures interface {
	GetOne(collectionName string, id string) (interface{}, error)
	CreateOne(collectionName string, data map[string]interface{}) (interface{}, error)
	UpdateOne(collectionName string, id string, data map[string]interface{}) (interface{}, error)
	DeleteOne(collectionName string, id string) (interface{}, error)
}

type Repo struct {
	Db *mongo.Database
}

func NewRepo() *Repo {
	db := db.NoSqlDbManager()
	return &Repo{db}
}

func (r *Repo) GetOne(collectionName string, id string) (interface{}, error) {

	var dataBsonM bson.M
	var data map[string]interface{}

	collection := r.Db.Collection(collectionName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&dataBsonM)
	if err != nil {
		return nil, err
	}

	err = helpers.JsonMarshaller(dataBsonM, &data)
	if err != nil {
		return nil, err
	}

	if data["encrypted"] == nil {
		return nil, errors.New("failed to decrypt")
	}

	decryptedString, err := helpers.AESGCMDecrypt(data["encrypted"].(string))
	if err != nil {
		return nil, err
	}

	var decryptedData interface{}
	err = json.Unmarshal([]byte(decryptedString), &decryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func (r *Repo) CreateOne(collectionName string, data map[string]interface{}) (interface{}, error) {

	var err error

	collection := r.Db.Collection(collectionName)

	id := primitive.NewObjectID()
	if data["_id"] != nil {
		id, err = primitive.ObjectIDFromHex(data["_id"].(string))
		if err != nil {
			return nil, err
		}
	}
	data["_id"] = id

	dataString, _ := json.Marshal(data)

	encryptedString, err := helpers.AESGCMEncrypt(string(dataString))
	if err != nil {
		return nil, err
	}

	encryptedData := map[string]interface{}{
		"_id":       data["_id"],
		"name":      data["name"],
		"encrypted": encryptedString,
	}

	result, err := collection.InsertOne(context.TODO(), encryptedData)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repo) UpdateOne(collectionName string, id string, data map[string]interface{}) (interface{}, error) {

	collection := r.Db.Collection(collectionName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	data["_id"] = objectId
	dataString, _ := json.Marshal(data)

	encryptedString, err := helpers.AESGCMEncrypt(string(dataString))
	if err != nil {
		return nil, err
	}

	encryptedData := map[string]interface{}{
		"encrypted": encryptedString,
	}

	toUpdate := bson.M{"$set": encryptedData}

	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objectId}, toUpdate)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repo) DeleteOne(collectionName string, id string) (interface{}, error) {

	collection := r.Db.Collection(collectionName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectId})

	if err != nil {
		return nil, err
	}

	return result, nil
}
