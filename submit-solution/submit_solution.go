package submitsolution

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/design/configs"
	"example.com/design/responses"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubmittedSolutionRequestBody struct {
	UserEmail   string `json:"user-name"`
	Solution    string `json:"solution-url"`
	ProblemName string `json:"problem-name"`
}

var validate = validator.New()

var unsettledCollection *mongo.Collection = configs.GetCollection(configs.DB, "unsettled")

func SubmitSolution() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var submittedSolutionRequestBody SubmittedSolutionRequestBody
		defer cancel()

		if err := json.NewDecoder(r.Body).Decode(&submittedSolutionRequestBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&submittedSolutionRequestBody); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		fmt.Println(submittedSolutionRequestBody)
		result, err := unsettledCollection.InsertOne(ctx, submittedSolutionRequestBody)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)

		response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)

	}
}
