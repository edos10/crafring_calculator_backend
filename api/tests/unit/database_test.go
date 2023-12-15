package main

import (
	"api_service/internal/databases"
	"api_service/tests/common"
	"database/sql"
	"sort"
	"testing"

	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/assert"
)

func IsNil(t any) bool {
	return t == nil
}

func CheckInputItems(t *testing.T, expected, actual []*databases.InputItem) {
	assert.True(t, IsNil(expected) == IsNil(actual))
	if IsNil(expected) {
		return
	}

	assert.Equal(t, len(expected), len(actual))

	sort.Slice(expected, func(i, j int) bool {
		return expected[i].ID < expected[j].ID
	})
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].ID < actual[j].ID
	})

	for i := range expected {
		assert.Equal(t, expected[i].ID, actual[i].ID)
		assert.Equal(t, expected[i].Quantity, actual[i].Quantity)
	}
}

func CheckItems(t *testing.T, expected, actual *databases.Item) {
	assert.True(t, IsNil(expected) == IsNil(actual))
	if IsNil(expected) {
		return
	}

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)

	assert.Equal(t, len(expected.Recipes), len(actual.Recipes))
	for i := range expected.Recipes {
		CheckRecipe(t, expected.Recipes[i], actual.Recipes[i])
	}
}

func CheckRecipe(t *testing.T, expected, actual *databases.Recipe) {
	assert.True(t, IsNil(expected) == IsNil(actual))
	if IsNil(expected) {
		return
	}

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ItemID, actual.ItemID)
	assert.Equal(t, expected.FactoryName, actual.FactoryName)
	assert.Equal(t, expected.ProductionFactory, actual.ProductionFactory)
	assert.Equal(t, expected.FactoryId, actual.FactoryId)
	assert.Equal(t, expected.BeltName, actual.BeltName)
	assert.Equal(t, expected.BeltQuantity, actual.BeltQuantity)

	assert.Equal(t, expected.InputItems, actual.InputItems)
}

func PrepareDatabase(t *testing.T) (*sql.DB, error) {
	db, err := sql.Open("ramsql", t.Name())
	if err != nil {
		return nil, err
	}
	if err = common.CreateTables(db); err != nil {
		return nil, err
	}
	if err = common.InsertData(db); err != nil {
		return nil, err
	}
	return db, nil
}

func PrepareDatabaseArrays(t *testing.T, recipes []*databases.Recipe, items []*databases.Item) (*sql.DB, error) {
	db, err := sql.Open("ramsql", t.Name())
	if err != nil {
		return nil, err
	}
	if err = common.CreateTables(db); err != nil {
		return nil, err
	}
	if err = common.InsertDataArrays(db, recipes, items); err != nil {
		return nil, err
	}
	return db, nil
}

func TestDatabase(t *testing.T) {
	recipes := []*databases.Recipe{
		&databases.Recipe{
			ID:                1,
			Name:              "Water Production",
			ItemID:            1,
			FactoryName:       "offshore pump",
			ProductionFactory: 72000,
			FactoryId:         1,
			BeltName:          "pipes",
			BeltQuantity:      17,
			InputItems:        nil,
		},
	}
	items := []*databases.Item{
		&databases.Item{
			ID:      1,
			Name:    "water",
			Recipes: recipes,
		},
	}
	db, err := PrepareDatabase(t)

	assert.NoError(t, err, "Failed to setup database")
	assert.NotNil(t, db, "Failed to setup database")

	sqlDb := databases.SqlDatabase{Connector: db}

	t.Run("TestItems", func(t *testing.T) {
		resp, err := sqlDb.GetItems()

		assert.NoErrorf(t, err, "GetItems should not fail")
		assert.NotNil(t, resp, "GetItems response should not be nil")

		assert.Equal(t, 1, len(resp), "GetItems should be size 1")
		CheckItems(t, items[0], resp[0])
	})

	t.Run("TestItem", func(t *testing.T) {
		resp, err := sqlDb.GetItem(1)

		assert.NoErrorf(t, err, "GetItem should not fail")
		assert.NotNil(t, resp, "GetItem response should not be nil")

		CheckItems(t, items[0], resp)
	})

	t.Run("TestRecipe", func(t *testing.T) {
		resp, err := sqlDb.GetRecipe(1)

		assert.NoErrorf(t, err, "GetRecipe should not fail")
		assert.NotNil(t, resp, "GetRecipe response should not be nil")

		CheckRecipe(t, recipes[0], resp)
	})
}

func TestDatabaseRecursive(t *testing.T) {
	recipes := []*databases.Recipe{
		&databases.Recipe{
			ID:                1,
			Name:              "Iron smelting",
			ItemID:            1,
			FactoryName:       "Furnace 5000",
			ProductionFactory: 60,
			FactoryId:         1,
			BeltName:          "pipes",
			BeltQuantity:      17,
			InputItems: []*databases.InputItem{
				&databases.InputItem{
					ID:       2,
					Quantity: 1,
				},
			},
		},
		&databases.Recipe{
			ID:                2,
			Name:              "Iron crafting",
			ItemID:            2,
			FactoryName:       "Crafting table",
			ProductionFactory: 5000000,
			FactoryId:         2,
			BeltName:          "pipes",
			BeltQuantity:      17,
			InputItems: []*databases.InputItem{
				&databases.InputItem{
					ID:       3,
					Quantity: 9,
				},
			},
		},
	}
	items := []*databases.Item{
		&databases.Item{
			ID:      1,
			Name:    "Iron Ingot",
			Recipes: []*databases.Recipe{recipes[0]},
		},
		&databases.Item{
			ID:      2,
			Name:    "Iron Ore",
			Recipes: []*databases.Recipe{recipes[1]},
		},
		&databases.Item{
			ID:      3,
			Name:    "Iron Nugget",
			Recipes: nil,
		},
	}
	db, err := PrepareDatabaseArrays(t, recipes, items)

	assert.NoError(t, err, "Failed to setup database")
	assert.NotNil(t, db, "Failed to setup database")

	sqlDb := databases.SqlDatabase{Connector: db}

	t.Run("TestItems", func(t *testing.T) {
		resp, err := sqlDb.GetItems()

		assert.NoErrorf(t, err, "GetItems should not fail")
		assert.NotNil(t, resp, "GetItems response should not be nil")

		assert.Equal(t, len(items), len(resp), "GetItems should be size 1")
		sort.Slice(resp, func(i, j int) bool {
			return resp[i].ID < resp[j].ID
		})
		for i := range items {
			CheckItems(t, items[i], resp[i])
		}
	})

	t.Run("TestItem", func(t *testing.T) {
		for i := range items {
			resp, err := sqlDb.GetItem(i + 1)

			assert.NoErrorf(t, err, "GetItem should not fail")
			assert.NotNil(t, resp, "GetItem response should not be nil")

			CheckItems(t, items[i], resp)
		}
	})

	t.Run("TestRecipe", func(t *testing.T) {
		for i := range recipes {
			resp, err := sqlDb.GetRecipe(i + 1)

			assert.NoErrorf(t, err, "GetRecipe should not fail")
			assert.NotNil(t, resp, "GetRecipe response should not be nil")

			CheckRecipe(t, recipes[i], resp)
		}
	})

	t.Run("TestRecursiveRecipe", func(t *testing.T) {
		resp, err := sqlDb.GetRecipeRecursive(1)

		assert.NoErrorf(t, err, "GetRecipe should not fail")
		assert.NotNil(t, resp, "GetRecipe response should not be nil")

		ind := 0
		for resp != nil {
			save := recipes[ind].InputItems

			recipes[ind].InputItems = nil
			CheckRecipe(t, recipes[ind], resp.ToRegular())
			recipes[ind].InputItems = save

			if len(resp.Children) == 0 {
				resp = nil
				continue
			}
			resp = resp.Children[0]
			ind += 1
			// FIXME(lexmach)
			if resp.BeltName == "FIXME" {
				resp = nil
				continue
			}
		}
	})
}
