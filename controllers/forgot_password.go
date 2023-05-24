package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/design/configs"
	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PasswordResetRequest struct {
	Email string `json:"email"`
}

func ForgotPassword(rw http.ResponseWriter, r *http.Request) {
	//we will have just user email as request body
	var req PasswordResetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// handle the error
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	Receiver := req.Email
	// var Receiver PasswordResetRequest
	// err := json.NewDecoder(r.Body).Decode(&Receiver)
	// fmt.Println(Receiver)
	// if err != nil {
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	RandomCode := CreateCode()
	SendMail(RandomCode, Receiver)

	if err != nil {
		rw.WriteHeader(http.StatusNotExtended)
		return
	}
	var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
	demoUser := FindUserByEmail(Receiver)
	demoId := demoUser.Id.Hex()
	fmt.Println(demoId)
	filter := bson.D{{"emailid", Receiver}}
	// fmt.Println(filter)
	update := bson.D{{"$set", bson.D{{"recoverycode", RandomCode},
		{"codesendingtime", time.Now()}}}}
	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(rw).Encode(response)
		return
	}

	// result, err := userCollection.UpdateOne(context.TODO(), filter, update)

	fmt.Println(result)
	// fmt.Println(demoUser.RecoveryCode)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(rw).Encode(response)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	response := responses.UserResponse{Status: http.StatusOK, Message: "Code Sent!", Data: map[string]interface{}{"data": Receiver}}
	json.NewEncoder(rw).Encode(response)

}
