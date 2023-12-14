package handlers

import (
	"api_service/internal/databases"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type RouteHandler struct {
	db databases.Database
}

// FIXME(lexmach): add log
func (handler *RouteHandler) getItems(w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		// log.Fatal("Trying to route getItems with nil router")
		w.WriteHeader(500)
		return
	}
	items, err := handler.db.GetItems()
	if err != nil {
		// log.Fatal(fmt.Sprintf("Trying to route getItems with error [%w] while getting items)", err))
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// FIXME(lexmach): add log
func (handler *RouteHandler) getRecipes(w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		w.WriteHeader(500)
		return
	}
	params := mux.Vars(r)
	itemID := params["item_id"]

	recipes, err := handler.db.GetRecipe(itemID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}
