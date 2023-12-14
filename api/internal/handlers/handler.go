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
	ID                int    `json:"recipe_id"`
	Name              string `json:"recipe_name"`
	ItemID            int    `json:"item_id"`
	FactoryName       string `json:"factory_name"`
	ProductionFactory int    `json:"production_factory"`
	FactoryId         int
	BeltName          string   `json:"belt_name"`
	BeltQuantity      int      `json:"belt_quantity"`
	Children          []Recipe `json:"children"`
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
	fmt.Println(itemID)
	recipes, err := getRecipeTree(itemID, db)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func getRecipeTree(itemId string, db *sql.DB) (*Recipe, error) {
	if itemId == "0" {
		return nil, nil
	}

	var recipe Recipe

	// читаем из recipes данные
	row := db.QueryRow("SELECT * FROM recipes WHERE item_id=$1", itemId)
	err := row.Scan(&recipe.ID, &recipe.Name, &recipe.ItemID, &recipe.FactoryId, &recipe.ProductionFactory)
	if err != nil {
		return nil, err
	}

	//
	rowForBelt := db.QueryRow("SELECT belt_id, quantity FROM recipe_belts WHERE recipe_id=$1", recipe.ID)

	var BeltId int
	err = rowForBelt.Scan(&BeltId, &recipe.BeltQuantity)
	if err != nil {
		return nil, err
	}

	rowForFactory := db.QueryRow("SELECT name FROM factories WHERE id=$1", recipe.FactoryId)

	err = rowForFactory.Scan(&recipe.FactoryName)
	if err != nil {
		return nil, err
	}

	//
	rowForBeltName := db.QueryRow("SELECT name FROM belts WHERE id=$1", BeltId)
	err = rowForBeltName.Scan(&recipe.BeltName)
	if err != nil {
		return nil, err
	}

	recipes_children := make([]Recipe, 0)

	rowsForChildRecipes, err := db.Query("SELECT child_id FROM recipes_ierarchy WHERE id=$1", recipe.ID)

	if err != nil {
		log.Fatal(err)
	}
	defer rowsForChildRecipes.Close()

	for rowsForChildRecipes.Next() {
		var child_id int
		err := rowsForChildRecipes.Scan(&child_id)
		if err != nil {
			return nil, err
		}

		child, err := getRecipeTree(fmt.Sprintf("%d", child_id), db)
		if err != nil {
			return nil, err
		}
		if child != nil {
			recipes_children = append(recipes_children, *child)
		}
	}

	recipe.Children = recipes_children
	return &recipe, nil
}

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/recipes/{item_id}", getRecipes).Methods("GET")
	return router
}
