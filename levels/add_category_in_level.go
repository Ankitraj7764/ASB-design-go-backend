package levels

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/design/models"
	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
)

type LevelWithCategory struct {
	LevelName     string            `json:"level-name"`
	ExtraCategory models.Categories `json:"extra-category"`
}

func AddCategoryInLevel() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var levelWithCategory LevelWithCategory

		if err := json.NewDecoder(r.Body).Decode(&levelWithCategory); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if validationErr := validate.Struct(&levelWithCategory); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		targetLevel := levelWithCategory.LevelName
		filter := bson.D{{"levelname", targetLevel}}
		var demoLevel models.Level
		levelCollection.FindOne(context.TODO(), filter).Decode(&demoLevel)
		demoCategoryList := demoLevel.CategoryList
		demoCategoryList = append(demoCategoryList, levelWithCategory.ExtraCategory)

		update := bson.D{{"$set", bson.D{{"categorylist", demoCategoryList}}}}
		result, err := levelCollection.UpdateOne(context.TODO(), filter, update)
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
