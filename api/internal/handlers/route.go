package handlers

import (
	"api_service/internal/databases"
	"encoding/json"
	"net/http"
	"strconv"
        "fmt"
	"github.com/gorilla/mux"
)

type RouteHandler struct {
	db databases.Database
}

// FIXME(lexmach): add log
func (handler *RouteHandler) getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if handler == nil {
		fmt.Fprintf(w, "Trying to route getItems with nil router")
		w.WriteHeader(500)
		return
	}
	items, err := handler.db.GetItems()
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Trying to route getItems with error [%w] while getting items)", err))
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// FIXME(lexmach): add log
func (handler *RouteHandler) getRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if handler == nil {
		w.WriteHeader(500)
		return
	}
	params := mux.Vars(r)
	itemID, err := strconv.Atoi(params["item_id"])

	if err != nil {
		w.WriteHeader(400)
		return
	}

	item, err := handler.db.GetItem(itemID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item.Recipes)
}
