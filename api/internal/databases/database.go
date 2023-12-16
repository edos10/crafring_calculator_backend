package databases

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type SqlDatabase struct {
	Connector *sql.DB
}

func GetSqlDatabse(dbHost, dbPort, dbUser, dbPassword, dbName string) (*SqlDatabase, error) {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &SqlDatabase{db}, err
}

// FIXME(lexmach): this is bad
func (db *SqlDatabase) GetRecipe(id RecipeID) (*Recipe, error) {
	if id == 0 {
		return nil, nil
	}

	recipe := &Recipe{}

	// читаем из recipes данные
	row := db.Connector.QueryRow("SELECT * FROM recipes WHERE item_id=$1", id)
	if row.Err() != nil && strings.Contains(row.Err().Error(), "no rows") {
		return nil, nil
	}
	err := row.Scan(&recipe.ID, &recipe.Name, &recipe.ItemID, &recipe.FactoryId, &recipe.ProductionFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in recipes with id %q: %w", id, err)
	}

	// FIXME(lexmach): rework
	rowForBelt := db.Connector.QueryRow("SELECT belt_id, quantity FROM recipe_belts WHERE recipe_id=$1", 1)
	var beltId int
	err = rowForBelt.Scan(&beltId, &recipe.BeltQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in recipe_belts with id %d: %w", recipe.ID, err)
	}

	rowForFactory := db.Connector.QueryRow("SELECT name FROM factories WHERE id=$1", recipe.FactoryId)

	err = rowForFactory.Scan(&recipe.FactoryName)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in factories with id %d: %w", recipe.FactoryId, err)
	}

	// FIXME(lexmach): rework
	rowForBeltName := db.Connector.QueryRow("SELECT name FROM belts WHERE id=$1", beltId)
	err = rowForBeltName.Scan(&recipe.BeltName)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in belts with id %d: %w", beltId, err)
	}

	rowsForChildRecipes, err := db.Connector.Query("SELECT item_id, quantity FROM recipes_input WHERE recipe_id=$1", recipe.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to query in recipes_input with id %d: %w", recipe.ID, err)
	}
	defer rowsForChildRecipes.Close()

	for rowsForChildRecipes.Next() {
		inputItem := &InputItem{}
		err := rowsForChildRecipes.Scan(&inputItem.ID, &inputItem.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to get childID in recipes_ierarchy: %w", err)
		}

		recipe.InputItems = append(recipe.InputItems, inputItem)
	}

	return recipe, nil
}

func (db *SqlDatabase) getItemRecipesId(id ItemID) (recipeIDs []RecipeID, err error) {
	if db == nil {
		return nil, fmt.Errorf("SqlDatabse is nil when fetching items")
	}
	if db.Connector == nil {
		return nil, fmt.Errorf("SqlDatabse connector is nil when fetching items")
	}

	rows, err := db.Connector.Query("SELECT id FROM recipes WHERE item_id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to select all items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var recipeID int
		err := rows.Scan(&recipeID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in getRecipesId: %w", err)
		}
		recipeIDs = append(recipeIDs, recipeID)
	}

	return recipeIDs, nil
}

func (db *SqlDatabase) GetItem(id ItemID) (item *Item, err error) {
	if db == nil {
		return nil, fmt.Errorf("SqlDatabse is nil when fetching items")
	}
	if db.Connector == nil {
		return nil, fmt.Errorf("SqlDatabse connector is nil when fetching items")
	}
	item = &Item{}

	row := db.Connector.QueryRow("SELECT * FROM items WHERE id=$1", id)
	err = row.Scan(&item.ID, &item.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in getItem %d: %w", id, err)
	}

	recipesIDs, err := db.getItemRecipesId(id)
	if err != nil {
		return nil, err
	}
	for _, recipeId := range recipesIDs {
		recipe, err := db.GetRecipe(recipeId)
		if err != nil {
			return nil, err
		}
		item.Recipes = append(item.Recipes, recipe)
	}

	return item, nil
}

func (db *SqlDatabase) GetItems() ([]*Item, error) {
	if db == nil {
		return nil, fmt.Errorf("SqlDatabse is nil when fetching items")
	}
	if db.Connector == nil {
		return nil, fmt.Errorf("SqlDatabse connector is nil when fetching items")
	}

	rows, err := db.Connector.Query("SELECT id FROM items")
	if err != nil {
		return nil, fmt.Errorf("failed to select all items: %w", err)
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var itemID int
		err := rows.Scan(&itemID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in getItems: %w", err)
		}

		item, err := db.GetItem(itemID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// FIXME(lexmach): need to count ProductionFactory correctly
func (db *SqlDatabase) GetRecipeRecursive(id RecipeID) (recipe *RecipeRecursive, err error) {
	if id == 0 {
		return nil, nil
	}

	recipeBase, err := db.GetRecipe(id)
	if err != nil {
		return nil, err
	}
	recipe = recipeBase.ToRecursive()

	for _, inputItem := range recipeBase.InputItems {
		inputItemDB, err := db.GetItem(inputItem.ID)
		if err != nil {
			return nil, err
		}

		// Bad case scenario, item has no recipes
		// this should be not possible, but who knows
		// TODO(lexmach): add log
		if len(inputItemDB.Recipes) == 0 {
			recipe.Children = append(recipe.Children, &RecipeRecursive{
				ID:                0,
				Name:              inputItemDB.Name,
				ItemID:            inputItem.ID,
				FactoryName:       "FIXME",
				ProductionFactory: 1,
				FactoryId:         0,
				BeltName:          "FIXME",
				BeltQuantity:      0,
				Children:          nil,
			})
			continue
		}
		// FIXME(lexmach): think of multiple recipes in inputItemDB
		inputItemRecipe, err := db.GetRecipeRecursive(inputItemDB.Recipes[0].ID)
		if err != nil {
			return nil, err
		}

		if inputItemRecipe != nil {
			recipe.Children = append(recipe.Children, inputItemRecipe)
		}
	}

	return recipe, nil
}

func CreateDatabaseConnect() (Database, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err := GetSqlDatabse(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
