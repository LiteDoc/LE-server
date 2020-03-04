package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type registerLock struct {
	Owner string
	Queue []string
}

func main() {
	godotenv.Load()
	port := 4000

	// Initialize register lock map
	registerLocks := map[string]registerLock{
		"1": registerLock{Owner: "none", Queue: []string{}},
		"2": registerLock{Owner: "none", Queue: []string{}},
		"3": registerLock{Owner: "none", Queue: []string{}},
	}

	// Initialize Router
	router := mux.NewRouter().StrictSlash(true)
	AddRoutes(router, registerLocks)

	// Start Server
	log.Print(fmt.Sprintf("Listening on port %d. \n", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))

}
