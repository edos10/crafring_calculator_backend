package handlers

import (
	"cdn_service/internal/databases"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type RouteHandler struct {
	db databases.Database
}

// FIXME(lexmach): add log
func (handler *RouteHandler) GetPath(w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		// log.Fatal("Trying to route getItems with nil router")
		w.WriteHeader(500)
		return
	}

	params := mux.Vars(r)
	itemID, err := strconv.Atoi(params["item_id"])

	if err != nil {
		w.WriteHeader(400)
		return
	}

	image, err := handler.db.GetPath(itemID)

	if err != nil {
		// log.Fatal(fmt.Sprintf("Trying to route getItems with error [%w] while getting items)", err))
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}
