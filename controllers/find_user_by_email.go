package controllers

import (
	"context"

	"example.com/design/configs"
	"example.com/design/models"
	"go.mongodb.org/mongo-driver/bson"
)

func FindUserByEmail(targetEmailId string) models.User {

	coll := configs.DB.Database("golangAPI").Collection("users")
	filter := bson.D{{"emailid", targetEmailId}}

	var result models.User
	coll.FindOne(context.TODO(), filter).Decode(&result)

	return result
}
