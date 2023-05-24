package categories

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/design/models"
	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
)

type CategoryWithChallenge struct {
	CategoryName   string           `json:"category-name"`
	ExtraChallenge models.Challange `json:"extra-challenge"`
}

func AddChallengeInCategory() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var categoryWithChallenge CategoryWithChallenge

		if err := json.NewDecoder(r.Body).Decode(&categoryWithChallenge); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		targetCategory := categoryWithChallenge.CategoryName
		filter := bson.D{{"categoryname", targetCategory}}
		var demoCategory models.Categories
		categoryCollection.FindOne(context.TODO(), filter).Decode(&demoCategory)
		demoProblemList := demoCategory.ProblemList
		demoProblemList = append(demoProblemList, categoryWithChallenge.ExtraChallenge)

		update := bson.D{{"$set", bson.D{{"problemlist", demoProblemList}}}}
		result, err := categoryCollection.UpdateOne(context.TODO(), filter, update)
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
