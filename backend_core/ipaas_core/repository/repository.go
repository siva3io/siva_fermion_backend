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
	"go.mongodb.org/mongo-driver/mongo/options"
)
/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License v3.0 as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License v3.0 for more details.
You should have received a copy of the GNU Lesser General Public License v3.0
along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/
type IpaasFeatures interface {
	GetOne(collectionName string, id string) (interface{}, error)
	CreateOne(collectionName string, data map[string]interface{}) (interface{}, error)
	UpdateOne(collectionName string, id string, data map[string]interface{}) (interface{}, error)
	DeleteOne(collectionName string, id string) (interface{}, error)
	GetMany(collectionName string, query map[string]interface{}) (interface{}, error)
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

	if collectionName == "ipaas_features" {
		encryptedData["app_code"] = data["app_code"]
		encryptedData["is_syncable"] = data["is_syncable"]
		encryptedData["display_name"] = data["display_name"]
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
		"name":      data["name"],
		"encrypted": encryptedString,
	}

	if collectionName == "ipaas_features" {
		encryptedData["app_code"] = data["app_code"]
		encryptedData["is_syncable"] = data["is_syncable"]
		encryptedData["display_name"] = data["display_name"]
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

func (r *Repo) GetMany(collectionName string, query map[string]interface{}) (interface{}, error) {

	var results []interface{}
	var bsonResults []bson.M

	collection := r.Db.Collection(collectionName)

	findOptions := options.Find().SetProjection(bson.M{"encrypted": 0})

	cur, err := collection.Find(context.TODO(), query, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem bson.M
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		bsonResults = append(bsonResults, elem)
	}

	err = helpers.JsonMarshaller(bsonResults, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
