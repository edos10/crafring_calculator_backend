package handlers

import (
	"cdn_service/internal/databases"
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
func (handler *RouteHandler) GetPath(w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		fmt.Fprintf(w, "Trying to route getPath unsuccess...")
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
		fmt.Fprintf(w, fmt.Sprintf("Trying to route getPath with error [%w] while getting path)", err))
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}
