package handlers

import (
	"api_service/internal/databases"

	"github.com/gorilla/mux"
)

func SetupRoutes() (*mux.Router, error) {
	db, err := databases.CreateDatabaseConnect()
	if err != nil {
		return nil, err
	}
	handler := &RouteHandler{db}

	router := mux.NewRouter()
	router.HandleFunc("/items", handler.getItems).Methods("GET")
	router.HandleFunc("/recipes/{item_id}", handler.getRecipes).Methods("GET")
	router.HandleFunc("/recipes/{item_id}", handler.getRecipes).Methods("GET")
	return router, nil
}
