package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type registerLock struct {
	owner string
	queue []string
}

func main() {
	godotenv.Load()
	port := 4000

	// Initialize register lock map
	registerLocks := map[string]registerLock{
		"1": registerLock{owner: "none", queue: []string{}},
		"2": registerLock{owner: "none", queue: []string{}},
		"3": registerLock{owner: "none", queue: []string{}},
	}

	// Initialize Router
	router := mux.NewRouter().StrictSlash(true)
	AddRoutes(router, registerLocks)

	// Start Server
	log.Print(fmt.Sprintf("Listening on port %d. \n", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))

}
