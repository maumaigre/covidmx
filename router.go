package main

import (
	"github.com/gorilla/mux"
)

// InitRouter creates a router and its routes
func InitRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/cases", getData).Methods("GET")
	router.HandleFunc("/stats", getStats).Methods("GET")

	router.HandleFunc("/forceFetch", forceFetch).Methods("POST")
	return router
}
