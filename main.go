package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"example.com/design/configs"
	"example.com/design/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//run database
	configs.ConnectDB()

	// //routes
	// routes.UserRoute(router)
	// log.Fatal(http.ListenAndServe(":8080", router))
	// Use the cors package to allow requests from all origins
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})
	routes.UserRoute(router)
	// Wrap the router with the cors middleware
	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
