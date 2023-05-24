package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/design/configs"
	"example.com/design/models"
	"example.com/design/responses"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser() http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("inside create user handler func")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		fmt.Println("request body has been decoded")

		//use the validator library to validate required fields
		// if validationErr := validate.Struct(&user); validationErr != nil {
		// 	rw.WriteHeader(http.StatusBadRequest)
		// 	response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
		// 	json.NewEncoder(rw).Encode(response)
		// 	return
		// }
		fmt.Println("validated")
		// emailIndex := mongo.IndexModel{
		// 	Keys:    bson.D{{"email-id", 1}},
		// 	Options: options.Index().SetUnique(true),
		// }
		// _, err := userCollection.Indexes().CreateOne(context.Background(), emailIndex)
		// if err != nil {
		// 	rw.WriteHeader(http.StatusBadRequest)
		// 	response := responses.UserResponse{Status: http.StatusBadRequest, Message: "email should be unique", Data: map[string]interface{}{"data": err.Error()}}
		// 	json.NewEncoder(rw).Encode(response)
		// 	return
		// }
		bytePassword := []byte(user.Password)
		hashedPassword := HashAndSalt(bytePassword)

		newUser := models.User{
			Id:   primitive.NewObjectID(),
			Name: user.Name,
			// Password: user.Password,
			Password: hashedPassword,
			// Location: user.Location,
			// Title:    user.Title,
			// Score:             user.Score,
			EmailId:           user.EmailId,
			ProfilePictureURL: user.ProfilePictureURL,
			RecoveryCode:      "-123",
			CodeSendingTime:   time.Now(),
		}
		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		// validate := validator.New()
		// validationError := validate.Struct(newUser)
		// if validationError != nil {
		// 	// email is not unique
		// 	fmt.Println("Email is not unique")
		// }

		rw.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
}
func GetAUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["userId"]
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

		if err != nil {
			// fmt.Println("user name is : ", user.Name)
			// fmt.Println("There is an error")
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		// fmt.Println(user.Name)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(rw).Encode(response)
	}
}

func EditAUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["userId"]
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		update := bson.M{
			"name":                user.Name,
			"id":                  user.Id,
			"email-id":            user.EmailId,
			"password":            user.Password,
			"profile-picture-url": user.ProfilePictureURL,
			"score":               user.Score,
			"solved-challenges":   user.SolvedChallenges,
			"recovery-code":       user.RecoveryCode,
			"code-sending-time":   user.CodeSendingTime,
		}
		result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//get updated user details
		var updatedUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}}
		json.NewEncoder(rw).Encode(response)
	}
}

func DeleteAUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["userId"]
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount < 1 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}}
		json.NewEncoder(rw).Encode(response)
	}
}
func GetAllUsers() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		var usersList []bson.M
		if err = cursor.All(ctx, &usersList); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		fmt.Println(usersList)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": usersList}}
		json.NewEncoder(rw).Encode(response)

	}
}
func GetLeaderboardDetails() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		// cursor, err := userCollection.Find(ctx, bson.M{})
		// if err != nil {
		// 	rw.WriteHeader(http.StatusInternalServerError)
		// 	response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		// 	json.NewEncoder(rw).Encode(response)
		// 	return
		// }
		// var usersList []bson.M
		// if err = cursor.All(ctx, &levels); err != nil {
		// 	rw.WriteHeader(http.StatusInternalServerError)
		// 	response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		// 	json.NewEncoder(rw).Encode(response)
		// 	return
		// }
		filter := bson.D{}
		opts := options.Find().SetSort(bson.D{{"score", -1}})
		cursor, err := userCollection.Find(context.TODO(), filter, opts)
		var results []models.User
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		for _, result := range results {
			res, _ := json.Marshal(result)
			fmt.Println(string(res))
		}
		// fmt.Println(result)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": results}}
		json.NewEncoder(rw).Encode(response)

	}
}
