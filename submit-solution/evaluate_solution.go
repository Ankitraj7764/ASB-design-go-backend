package submitsolution

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/design/configs"
	"example.com/design/controllers"
	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

//what will be the request body containing?
/*
create a queue/stack/array like design in the database
where every element(pendingProblemRequestBody) in the queue will contain these 3 things
	{
		username
		problem name
		solutionURL
	}

	now admin will be able to see all these pendingProblemRequestBody
	they will have two options:
		- Accept it
		- Reject it

	If accepted it will check that for that user if the problem is already solved
		- If it's already solved we will not increase the score
		- Else we will increase the score and put the problem in solved challenge array for that user

func EvaluateSolution() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	}
}
*/

// var validate = validator.New()

// type PendingProblemRequestBody struct {
// 	UserName    string `json:"user-name"`
// 	ProblemName string `json:"Problem-name"`
// 	SolutionURL string `json:"solution-url"`
// }

func AcceptSolution() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var pendingProblemRequestBody SubmittedSolutionRequestBody
		defer cancel()
		fmt.Println("inside accecpt-solution function")
		err := json.NewDecoder(r.Body).Decode(&pendingProblemRequestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		fmt.Println("request body decoded")
		validationErr := validate.Struct(&pendingProblemRequestBody)
		if validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "validation error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		fmt.Println("request body validated")

		// extract user using username
		targetUser := controllers.FindUserByEmail(pendingProblemRequestBody.UserEmail)
		fmt.Println("target user found")
		// if targetUser == nil{
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error : can't find problem with given name", Data: map[string]interface{}{}}
		// 	json.NewEncoder(w).Encode(response)
		// 	return
		// }
		solvedChallengesArr := targetUser.SolvedChallenges

		fmt.Println("found old array of solved problems")
		// extract problem using problem name
		targetChallenge, err := FindChallengeByName(pendingProblemRequestBody.ProblemName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		fmt.Println("target challenge found")
		for i := 0; i < len(solvedChallengesArr); i++ {
			if solvedChallengesArr[i] == targetChallenge {
				w.WriteHeader(http.StatusAlreadyReported)
				response := responses.UserResponse{Status: http.StatusAlreadyReported, Message: "Problem has already been solved", Data: map[string]interface{}{}}
				json.NewEncoder(w).Encode(response)
				return
			}

		}
		previousScore := targetUser.Score
		challengeScore := targetChallenge.Score
		var newScore int = previousScore + challengeScore

		fmt.Println("challenge was not there previously")
		//update array
		solvedChallengesArr = append(solvedChallengesArr, targetChallenge)
		fmt.Println("challenge added to solved challenge array")
		//update user
		filter := bson.D{{"emailid", pendingProblemRequestBody.UserEmail}}

		update := bson.D{{"$set", bson.D{{"solvedchallenges", solvedChallengesArr}, {"score", newScore}}}}
		var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
		result, err2 := userCollection.UpdateOne(context.TODO(), filter, update)
		if err2 != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err2.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		fmt.Println("user details updated")

		res, err := unsettledCollection.DeleteOne(ctx, bson.M{"problemname": targetChallenge.ChallengeName})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("deleted %v documents\n", res.DeletedCount)

		w.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(w).Encode(response)

	}
}
