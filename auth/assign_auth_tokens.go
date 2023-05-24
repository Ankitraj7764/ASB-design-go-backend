package auth

import (
	"context"
	"fmt"
	"time"

	"example.com/design/configs"
	"example.com/design/controllers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// E represents a BSON element for a D. It is usually used inside a D
type E struct {
	Key   string
	Value interface{}
}

// AssignAuthToken takes email string as argument and
// creates access token and refresh token
// and append refresh token the list-of-refresh token
// for that corresponding user

func AssignAuthTokens(email string) (string, string, error) {

	//generate access token

	newAccessToken, err := CreateNewJWTAccessToken(email)
	if err != nil {
		return "", "", err
	}
	// newAccessTokenByte := []byte(newAccessToken)
	fmt.Println(newAccessToken)
	//generate refresh token

	newRefreshToken := CreateNewReferenceToken(6)
	newRefreshTokenByte := []byte(newRefreshToken)

	//encrypt access token

	// newEncryptedAccessToken := controllers.HashAndSalt(newAccessTokenByte)

	//encrypt refresh token

	newEncryptedRefreshToken := controllers.HashAndSalt(newRefreshTokenByte)

	//clear the refresh tokens array

	var newRefreshTokenArray []string

	//add new encrypted refresh token to the array

	newRefreshTokenArray = append(newRefreshTokenArray, newEncryptedRefreshToken)

	filter := bson.D{primitive.E{Key: "emailid", Value: email}}
	update := bson.D{{"$set", bson.D{primitive.E{Key: "refreshtokens", Value: newRefreshTokenArray}, primitive.E{Key: "refreshtokensendingtime", Value: time.Now()}}}}

	_, err2 := userCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		return "", "", err2
	}

	return newAccessToken, newRefreshToken, nil
}
