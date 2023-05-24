package modules

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"example.com/design/configs"
	"example.com/design/models"
	"example.com/design/responses"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

var moduleCollection *mongo.Collection = configs.GetCollection(configs.DB, "modules")
var validate = validator.New()

func CreateModule() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var module models.Module
		defer cancel()

		if err := json.NewDecoder(r.Body).Decode(&module); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&module); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newModule := models.Module{
			// Id:         primitive.NewObjectID(),
			ModuleName: module.ModuleName,
			LevelList:  module.LevelList,
			// Password: user.Password,

			// Location: user.Location,
			// Title:    user.Title,
			// Score:             user.Score,
			// EmailId:           user.EmailId,
			// ProfilePictureURL: user.ProfilePictureURL,
			// RecoveryCode:      "-123",
			// CodeSendingTime:   time.Now(),

		}
		result, err := moduleCollection.InsertOne(ctx, newModule)
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
