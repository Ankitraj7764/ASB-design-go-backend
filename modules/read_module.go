package modules

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/design/models"
	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
)

type ModuleRequestBody struct {
	ModuleName string `json:"module-name"`
}

func FindModule() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var moduleReqBody ModuleRequestBody

		if err := json.NewDecoder(r.Body).Decode(&moduleReqBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if validationErr := validate.Struct(&moduleReqBody); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		// coll := configs.DB.Database("golangAPI").Collection("modules")
		// fmt.Println(moduleReqBody.ModuleName)
		targetModule := moduleReqBody.ModuleName
		filter := bson.D{{"modulename", targetModule}}

		var result models.Module
		err := moduleCollection.FindOne(context.TODO(), filter).Decode(&result)
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
