package levels

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/design/models"
	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
)

type LevelRequestBody struct {
	LevelName string `json:"level-name"`
}

func FindLevel() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var levelRequestBody LevelRequestBody

		if err := json.NewDecoder(r.Body).Decode(&levelRequestBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if validationErr := validate.Struct(&levelRequestBody); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		// coll := configs.DB.Database("golangAPI").Collection("modules")
		// fmt.Println(moduleReqBody.ModuleName)
		targetLevel := levelRequestBody.LevelName
		filter := bson.D{{"levelname", targetLevel}}

		var result models.Level
		err := levelCollection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusFound)
		response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
}
