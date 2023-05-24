package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandlePage(writer http.ResponseWriter, request *http.Request) {
	_, err := GenerateJWT("random email")
	if err != nil {
		log.Fatalln("Error generating JWT", err)
	}

	writer.Header().Set("Token", "%v")
	type_ := "application/json"
	writer.Header().Set("Content-Type", type_)
	var message Message
	err = json.NewDecoder(request.Body).Decode(&message)
	if err != nil {
		return
	}
	err = json.NewEncoder(writer).Encode(message)
	if err != nil {
		return
	}
}
