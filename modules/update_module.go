package modules

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/design/models"
	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
)

type ModuleWithLevel struct {
	ModuleName string       `json:"module-name"`
	ExtraLevel models.Level `json:"extra-level"`
}

func AddLevelInModule() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var moduleWithLevel ModuleWithLevel

		if err := json.NewDecoder(r.Body).Decode(&moduleWithLevel); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if validationErr := validate.Struct(&moduleWithLevel); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		targetModule := moduleWithLevel.ModuleName
		filter := bson.D{{"modulename", targetModule}}
		var demoModule models.Module
		moduleCollection.FindOne(context.TODO(), filter).Decode(&demoModule)
		demoLevelList := demoModule.LevelList
		demoLevelList = append(demoLevelList, moduleWithLevel.ExtraLevel)

		update := bson.D{{"$set", bson.D{{"levellist", demoLevelList}}}}
		result, err := moduleCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		rw.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
}
