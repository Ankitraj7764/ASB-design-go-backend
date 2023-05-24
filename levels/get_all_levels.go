package levels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllLevels() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		cursor, err := levelCollection.Find(ctx, bson.M{})
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		var levels []bson.M
		if err = cursor.All(ctx, &levels); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		fmt.Println(levels)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": levels}}
		json.NewEncoder(rw).Encode(response)

	}
}
