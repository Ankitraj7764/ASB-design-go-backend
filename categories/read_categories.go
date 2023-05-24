package categories

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/design/models"
	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
)

type CategoryRequestBody struct {
	CategoryName string `json:"category-name"`
}

func FindCategory() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var categoryRequestBody CategoryRequestBody

		if err := json.NewDecoder(r.Body).Decode(&categoryRequestBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if validationErr := validate.Struct(&categoryRequestBody); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		// coll := configs.DB.Database("golangAPI").Collection("modules")
		// fmt.Println(moduleReqBody.ModuleName)
		targetCategory := categoryRequestBody.CategoryName
		filter := bson.D{{"categoryname", targetCategory}}

		var result models.Level
		err := categoryCollection.FindOne(context.TODO(), filter).Decode(&result)
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
