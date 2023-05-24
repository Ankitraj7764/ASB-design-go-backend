package validators

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"example.com/design/responses"
)

type RequestBody struct {
	Email        string `json:"email"`
	RefreshToken string `json:"refresh-token"`
}
type ResponseBody struct {
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
}

func UseRefreshToken() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()
		var requestBody RequestBody
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "unable to decode the request body", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		newAccessToken, newRefreshToken, err := GetNewAccessTokenByRefreshToken(requestBody.RefreshToken, requestBody.Email)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusExpectationFailed, Message: "unable to generate new tokens", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		var result ResponseBody
		result.RefreshToken = newRefreshToken
		result.AccessToken = newAccessToken
		rw.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)

	}
}
