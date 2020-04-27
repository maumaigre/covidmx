package main

import (
	"github.com/gorilla/mux"
)

// InitRouter creates a router and its routes
func InitRouter() *mux.Router {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api").Subrouter()
	subRouter.HandleFunc("/", getMain).Methods("GET")
	subRouter.HandleFunc("/cases", getData).Methods("GET")
	subRouter.HandleFunc("/stats", getStats).Methods("GET")

	subRouter.HandleFunc("/stateStats", getStateStats).Methods("GET")

	subRouter.HandleFunc("/dailyNewStats", getDailyNewStats).Methods("GET")

	subRouter.HandleFunc("/forceFetch", forceFetch).Methods("POST")
	return router
}
