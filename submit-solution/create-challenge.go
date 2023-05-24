package submitsolution

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"example.com/design/configs"
	"example.com/design/models"
	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/mongo"
)

var challengeCollection *mongo.Collection = configs.GetCollection(configs.DB, "challenges")

func CreateChallenge() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var challengeRequestBody models.Challange
		defer cancel()
		if err := json.NewDecoder(r.Body).Decode(&challengeRequestBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			resonse := responses.UserResponse{Status: http.StatusBadRequest, Message: "Can't decode the request body of HTTP requst!", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(resonse)
			return
		}
		if validationErr := validate.Struct(challengeRequestBody); validationErr != nil {
			rw.WriteHeader(http.StatusNonAuthoritativeInfo)
			response := responses.UserResponse{Status: http.StatusNonAuthoritativeInfo, Message: "validation Error!", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		result, err := challengeCollection.InsertOne(ctx, challengeRequestBody)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "Unable to insert new Challenge!", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		rw.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{Status: http.StatusCreated, Message: "Success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)

	}
}
