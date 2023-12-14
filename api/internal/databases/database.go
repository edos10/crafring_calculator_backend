package databases

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type SqlDatabase struct {
	connector *sql.DB
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

func (db *SqlDatabase) GetItems() ([]*Item, error) {
	if db == nil {
		return nil, fmt.Errorf("SqlDatabse is nil when fetching items")
	}
	if db.connector == nil {
		return nil, fmt.Errorf("SqlDatabse connector is nil when fetching items")
	}

	rows, err := db.connector.Query("SELECT * FROM items")
	if err != nil {
		return nil, fmt.Errorf("failed to select all items: %w", err)
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item := &Item{}
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// FIXME(lexmach): logic refactor
func (db *SqlDatabase) GetRecipe(id string) (*Recipe, error) {
	if id == "0" {
		return nil, nil
	}

	recipe := &Recipe{}

	// читаем из recipes данные
	// FIXME(lexmach): sql injections fix
	row := db.connector.QueryRow("SELECT * FROM recipes WHERE item_id=$1", id)
	err := row.Scan(&recipe.ID, &recipe.Name, &recipe.ItemID, &recipe.FactoryId, &recipe.ProductionFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in recipes with id %q: %w", id, err)
	}

	//
	rowForBelt := db.connector.QueryRow("SELECT belt_id, quantity FROM recipe_belts WHERE recipe_id=$1", recipe.ID)

	var beltId int
	err = rowForBelt.Scan(&beltId, &recipe.BeltQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in recipe_belts with id %d: %w", recipe.ID, err)
	}

	rowForFactory := db.connector.QueryRow("SELECT name FROM factories WHERE id=$1", recipe.FactoryId)

	err = rowForFactory.Scan(&recipe.FactoryName)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in factories with id %d: %w", recipe.FactoryId, err)
	}

	//
	rowForBeltName := db.connector.QueryRow("SELECT name FROM belts WHERE id=$1", beltId)
	err = rowForBeltName.Scan(&recipe.BeltName)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in belts with id %d: %w", beltId, err)
	}

	recipe.Children = make([]*Recipe, 0)

	rowsForChildRecipes, err := db.connector.Query("SELECT child_id FROM recipes_ierarchy WHERE id=$1", recipe.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to query in recipes_ierarchy with id %d: %w", recipe.ID, err)
	}
	defer rowsForChildRecipes.Close()

	for rowsForChildRecipes.Next() {
		var childID int
		err := rowsForChildRecipes.Scan(&childID)
		if err != nil {
			return nil, fmt.Errorf("failed to get childID in recipes_ierarchy: %w", err)
		}

		child, err := db.GetRecipe(fmt.Sprintf("%d", childID))
		if err != nil {
			return nil, err
		}

		if child != nil {
			recipe.Children = append(recipe.Children, child)
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
