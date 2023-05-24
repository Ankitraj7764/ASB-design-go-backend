package submitsolution

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"example.com/design/responses"
	"gopkg.in/mgo.v2/bson"
)

func GetAllChallenges() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		cursor, err := challengeCollection.Find(ctx, bson.M{})
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		var challenges []bson.M
		if err = cursor.All(ctx, &challenges); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		// fmt.Println(challenges)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": challenges}}
		json.NewEncoder(rw).Encode(response)

	}
}
