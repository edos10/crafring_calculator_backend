package handlers

import (
	"api_service/internal/databases"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Item struct {
	ID   int    `json:"item_id"`
	Name string `json:"item_name"`
}

type Recipe struct {
	ID       int      `json:"recipe_id"`
	Name     string   `json:"recipe_name"`
	ItemID   int      `json:"item_id"`
	Children []Recipe `json:"children"`
}

func getItems(w http.ResponseWriter, r *http.Request) {
	db, err := databases.CreateDatabaseConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func getRecipes(w http.ResponseWriter, r *http.Request) {
	db, err := databases.CreateDatabaseConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	itemID := params["item_id"]

	recipes, err := getRecipeTree(itemID, db)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func getRecipeTree(itemID string, db *sql.DB) ([]Recipe, error) {
	rows, err := db.Query("SELECT * FROM recipes WHERE item_id=$1", itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.ItemID)
		if err != nil {
			return nil, err
		}

		children, err := getRecipeTree(fmt.Sprintf("%d", recipe.ID), db)
		if err != nil {
			return nil, err
		}
		recipe.Children = children

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/recipes/{item_id}", getRecipes).Methods("GET")
	return router
}
