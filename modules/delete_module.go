package modules

// package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/design/responses"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteModule() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// params := mux.Vars(r)
		// userId := params["userId"]
		defer cancel()
		var moduleRequestBody ModuleRequestBody
		if err := json.NewDecoder(r.Body).Decode(&moduleRequestBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if validationErr := validate.Struct(&moduleRequestBody); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		// objId, _ := primitive.ObjectIDFromHex(userId)

		targetModuleName := moduleRequestBody.ModuleName
		fmt.Println(targetModuleName)
		result, err := moduleCollection.DeleteOne(ctx, bson.M{"modulename": targetModuleName})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount < 1 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Module with that name not found!"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}}
		json.NewEncoder(rw).Encode(response)
	}
}

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"

// 	"example.com/design/responses"
// 	"go.mongodb.org/mongo-driver/bson"
// )

// func DeleteModule() http.HandlerFunc {
// 	return func(rw http.ResponseWriter, r *http.Request) {
// 		var moduleReqBody ModuleRequestBody

// 		if err := json.NewDecoder(r.Body).Decode(&moduleReqBody); err != nil {
// 			rw.WriteHeader(http.StatusBadRequest)
// 			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
// 			json.NewEncoder(rw).Encode(response)
// 			return
// 		}

// 		if validationErr := validate.Struct(&moduleReqBody); validationErr != nil {
// 			rw.WriteHeader(http.StatusBadRequest)
// 			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
// 			json.NewEncoder(rw).Encode(response)
// 			return
// 		}
// 		// coll := configs.DB.Database("golangAPI").Collection("modules")
// 		// fmt.Println(moduleReqBody.ModuleName)
// 		targetModule := moduleReqBody.ModuleName
// 		// session, err := mgo.Dial("localhost")
// 		// if err != nil {
// 		// 	fmt.Printf("dial fail %v\n", err)
// 		// 	os.Exit(1)
// 		// }
// 		// defer session.Close()
// 		// filter := bson.D{{"modulename", targetModule}}

// 		// var result models.Module
// 		// err := moduleCollection.FindOne(context.TODO(), filter).Decode(&result)
// 		// newModuleCollection := session.DB("golangAPI").C("modules")
// 		err := moduleCollection.Remove(bson.M{"modulename": targetModule})

// 		if err != nil {
// 			rw.WriteHeader(http.StatusBadRequest)
// 			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
// 			json.NewEncoder(rw).Encode(response)
// 			return
// 		}

// 		rw.WriteHeader(http.StatusFound)
// 		response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": targetModule}}
// 		json.NewEncoder(rw).Encode(response)
// 	}
// }
