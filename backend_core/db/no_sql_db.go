package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var noSqlDb *mongo.Database

func NoSqlInit() {

	MONGO_DB_HOST := os.Getenv("MONGO_DB_HOST")
	MONGO_DB_USER := os.Getenv("MONGO_DB_USER")
	MONGO_DB_PASS := os.Getenv("MONGO_DB_PASS")
	MONGO_DB_NAME := os.Getenv("MONGO_DB_NAME")
	MONGO_DB_PORT := os.Getenv("MONGO_DB_PORT")
	MONGO_DB_SRV := os.Getenv("MONGO_DB_SRV")

	var MONGO_DB_USERNAME_AND_PASSWORD string

	if MONGO_DB_USER != "" && MONGO_DB_PASS != "" {
		MONGO_DB_USERNAME_AND_PASSWORD = fmt.Sprintf("%v:%v@", MONGO_DB_USER, MONGO_DB_PASS)
	}

	MONGO_DB_HOST_AND_PORT := MONGO_DB_HOST
	if MONGO_DB_PORT != "" {
		MONGO_DB_HOST_AND_PORT = fmt.Sprintf("%v:%v", MONGO_DB_HOST, MONGO_DB_PORT)
	}

	var MONGO_DB_SRV_VALUE string
	if MONGO_DB_SRV != "false" {
		MONGO_DB_SRV_VALUE = "+srv"
	}

	connectionString := fmt.Sprintf("mongodb%v://%v%v", MONGO_DB_SRV_VALUE, MONGO_DB_USERNAME_AND_PASSWORD, MONGO_DB_HOST_AND_PORT)

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		panic("failed to connect no_sql database")
	}
	if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic("failed to connect no_sql database")
	}

	noSqlDb = mongoClient.Database(MONGO_DB_NAME)

}

func NoSqlDbManager() *mongo.Database {
	return noSqlDb
}
