package levels

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

var levelCollection *mongo.Collection = configs.GetCollection(configs.DB, "levels")
var validate = validator.New()

func CreateLevel() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var level models.Level
		defer cancel()

		if err := json.NewDecoder(r.Body).Decode(&level); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&level); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newLevel := models.Level{
			// Id:         primitive.NewObjectID(),
			// ModuleName: module.ModuleName,
			// LevelList:  module.LevelList,
			LevelName:    level.LevelName,
			CategoryList: level.CategoryList,
			// Password: user.Password,

			// Location: user.Location,
			// Title:    user.Title,
			// Score:             user.Score,
			// EmailId:           user.EmailId,
			// ProfilePictureURL: user.ProfilePictureURL,
			// RecoveryCode:      "-123",
			// CodeSendingTime:   time.Now(),

		}
		result, err := levelCollection.InsertOne(ctx, newLevel)
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
