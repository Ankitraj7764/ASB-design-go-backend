package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/design/responses"
	"example.com/design/validators"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func LogOut(rw http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("validator")
	err := validators.ValidateAccessToken(bearer)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(rw).Encode(response)
		return
	}
	var newRefreshTokenArray []string
	// fmt.Println("bearer ", bearer)
	//add new encrypted refresh token to the array
	// var email string
	// extract email
	// for key, val := range bearer {
	// 	if key == "emailid" {
	// 		email = val
	// 		break
	// 	}
	// }
	email, err := extractClaims(bearer)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error while extracting claims", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(rw).Encode(response)
		return
	}
	// fmt.Println("at line 36")
	filter := bson.D{primitive.E{Key: "emailid", Value: email}}
	update := bson.D{{"$set", bson.D{primitive.E{Key: "refreshtokens", Value: newRefreshTokenArray}}}}

	// fmt.Println("at line 40")
	_, err2 := userCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error while updating data", Data: map[string]interface{}{"data": err2.Error()}}
		json.NewEncoder(rw).Encode(response)
		return
	}

}

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

var sampleSecretKey = []byte("SecretYouShouldHide")

func extractClaims(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("there's an error with the signing method")
		}
		return sampleSecretKey, nil
	})
	if err != nil {
		return "Error Parsing Token: ", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email := claims["email"].(string)
		return email, nil
	}

	var result string = "unable to extract claims"
	return result, nil
}
