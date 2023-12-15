package common

import (
	"api_service/internal/databases"
	"database/sql"
	"fmt"

	"github.com/tanimutomo/sqlfile"
)

func CreateTables(db *sql.DB) error {
	s := sqlfile.New()
	err := s.File("../common/data/create_tables.sql")
	if err != nil {
		return err
	}
	_, err = s.Exec(db)
	if err != nil {
		return err
	}
	return nil
}

func InsertData(db *sql.DB) error {
	s := sqlfile.New()
	err := s.File("../common/data/insert_test_data.sql")
	if err != nil {
		return err
	}
	_, err = s.Exec(db)
	if err != nil {
		return err
	}
	return nil
}

func InsertDataArrays(db *sql.DB, recipes []*databases.Recipe, items []*databases.Item) error {
	_, err := db.Exec("INSERT INTO belts(id, name) VALUES (1, \"pipes\");")
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO recipe_belts (id, recipe_id, belt_id, quantity) VALUES (1, 1, 1, 17);")
	if err != nil {
		return err
	}

	factories := make(map[int]string)
	item_map := make(map[int]string)
	for _, recipe := range recipes {
		if val, ok := factories[recipe.FactoryId]; ok && val != recipe.FactoryName {
			return fmt.Errorf("duplicate factoryID with different names; (%d, %q), (%d, %q)", recipe.FactoryId, recipe.FactoryName, recipe.FactoryId, val)
		}
		factories[recipe.FactoryId] = recipe.FactoryName

		_, err = db.Exec("INSERT INTO recipes (id, name, item_id, factory_id, production_rate_per_second) VALUES ($1, $2, $3, $4, $5)",
			recipe.ID, recipe.Name, recipe.ItemID, recipe.FactoryId, recipe.ProductionFactory)

		if err != nil {
			return err
		}

		for _, input := range recipe.InputItems {
			_, err = db.Exec("INSERT INTO recipes_input (recipe_id, item_id, item_quantity) VALUES ($1, $2, $3)", recipe.ID, input.ID, input.Quantity)

			if err != nil {
				return err
			}
		}
	}
	for _, item := range items {
		if val, ok := item_map[item.ID]; ok && val != item.Name {
			return fmt.Errorf("duplicate itemID with different names; (%d, %q), (%d, %q)", item.ID, item.Name, item.ID, item.Name)
		}
		item_map[item.ID] = item.Name
	}

	for id, name := range factories {
		_, err := db.Exec("INSERT INTO factories (id, name) VALUES ($1, $2)", id, name)
		if err != nil {
			return err
		}
	}
	for id, name := range item_map {
		_, err := db.Exec("INSERT INTO items (id, name) VALUES ($1, $2)", id, name)
		if err != nil {
			return err
		}
	}
	return nil
}
