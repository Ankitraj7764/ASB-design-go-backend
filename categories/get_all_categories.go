package categories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllCategories() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		cursor, err := categoryCollection.Find(ctx, bson.M{})
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		var categories []bson.M
		if err = cursor.All(ctx, &categories); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		fmt.Println(categories)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": categories}}
		json.NewEncoder(rw).Encode(response)

	}
}
