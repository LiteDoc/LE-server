package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type jSONResponse struct {
	Status string
	Code   int
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// readLockHandler : process lock request from proxy server
func readLockHandler(w http.ResponseWriter, r *http.Request, registerLocks map[string]registerLock) {
	// Get query params
	var q = r.URL.Query()
	var clientName string = q.Get("name")
	var clientRegisterID string = q.Get("registerID")

	if registerLocks[clientRegisterID].Owner == clientName {
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
	return
}

// readLocksHandler : return information of register locks
func readLocksHandler(w http.ResponseWriter, r *http.Request, registerLocks map[string]registerLock) {
	log.Print("Register locks read.\n")
	json.NewEncoder(w).Encode(registerLocks)
	return
}

// lockHandler : process lock request from proxy server
func lockHandler(w http.ResponseWriter, r *http.Request, registerLocks map[string]registerLock) {
	// Get query params
	var q = r.URL.Query()
	var clientName string = q.Get("name")
	var clientRegisterID string = q.Get("registerID")

	var selectedLock = registerLocks[clientRegisterID]
	if selectedLock.Owner == "none" {
		selectedLock.Owner = clientName
		registerLocks[clientRegisterID] = selectedLock
		log.Print(registerLocks)
		json.NewEncoder(w).Encode(jSONResponse{Status: "Lock Success", Code: 200})
		return
	}
	log.Print(fmt.Sprintf("Failed Lock attempt by %s\n", clientName))
	json.NewEncoder(w).Encode(jSONResponse{Status: "Lock Failed", Code: 401})
	return
}

// unlockHandler : process unlock request from proxy server
func unlockHandler(w http.ResponseWriter, r *http.Request, registerLocks map[string]registerLock) {
	// Get query params
	var q = r.URL.Query()
	var clientName string = q.Get("name")
	var clientRegisterID string = q.Get("registerID")

	var selectedLock = registerLocks[clientRegisterID]
	if selectedLock.Owner == clientName {
		selectedLock.Owner = "none"
		registerLocks[clientRegisterID] = selectedLock
		log.Print(registerLocks)
		json.NewEncoder(w).Encode(jSONResponse{Status: "Unlock Success", Code: 200})
		return
	}
	log.Print(fmt.Sprintf("Failed Unlock attempt by %s\n", clientName))
	json.NewEncoder(w).Encode(jSONResponse{Status: "Unlock Failed", Code: 401})
	return
}

// AddRoutes : attach route handlers to router
func AddRoutes(router *mux.Router, registerLocks map[string]registerLock) {
	router.HandleFunc("/readLock", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		readLockHandler(w, r, registerLocks)
	})
	router.HandleFunc("/readLocks", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		readLocksHandler(w, r, registerLocks)
	})
	router.HandleFunc("/lock", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		lockHandler(w, r, registerLocks)
	})
	router.HandleFunc("/unlock", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		unlockHandler(w, r, registerLocks)
	})

}
