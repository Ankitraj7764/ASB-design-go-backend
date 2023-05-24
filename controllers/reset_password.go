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

type NewData struct {
	Email       string `json:"email"`
	Code        string `json:"recovery-code"`
	NewPassword string `json:"password"`
}

func ResetPassword(rw http.ResponseWriter, r *http.Request) {
	//we will get  code and password as the body
	var newData NewData
	err := json.NewDecoder(r.Body).Decode(&newData)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	f := CompareCodes(newData.Code, newData.Email)
	if !f {
		rw.WriteHeader(http.StatusUnavailableForLegalReasons)
		return
	}
	currentTime := GetCurrentTime()
	demoUser := FindUserByEmail(newData.Email)
	sentTime := demoUser.CodeSendingTime
	elapsed := currentTime.Sub(sentTime)
	if elapsed > 15*time.Minute {
		rw.WriteHeader(http.StatusGatewayTimeout)
		return
	}
	// reset password
	newPassword := newData.NewPassword
	newBytePassword := []byte(newPassword)
	newHashedPassword := HashAndSalt(newBytePassword)
	var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
	// demoUser := FindUserByEmail(newData.Email)
	// demoId := demoUser.Id.Hex()
	// fmt.Println(demoId)
	filter := bson.D{{"emailid", newData.Email}}
	// fmt.Println(filter)
	update := bson.D{{"$set", bson.D{{"password", newHashedPassword}}}}
	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	// display the number of documents updated
	fmt.Println("Number of documents updated:", result.ModifiedCount)

	rw.WriteHeader(http.StatusCreated)
	response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": demoUser}}
	json.NewEncoder(rw).Encode(response)

}

func CompareCodes(Code, EnteredEmail string) bool {

	demoUser := FindUserByEmail(EnteredEmail)
	fmt.Println(Code)
	fmt.Println(demoUser.RecoveryCode)
	fmt.Println(demoUser.RecoveryCode, " ", Code)
	if demoUser.RecoveryCode == Code {
		return true
	}

	return false

}
