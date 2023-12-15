package handlers

import (
	"cdn_service/internal/databases"

	"github.com/gorilla/mux"
)

func SetupRoutes() (*mux.Router, error) {
	db, err := databases.CreateDatabaseConnect()
	if err != nil {
		return nil, err
	}
	handler := &RouteHandler{db}

	router := mux.NewRouter()
	router.HandleFunc("/path/{item_id}", handler.GetPath).Methods("GET")
	return router, nil
}
