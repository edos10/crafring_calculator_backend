package databases

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
	BeltName          string    `json:"belt_name"`
	BeltQuantity      int       `json:"belt_quantity"`
	Children          []*Recipe `json:"children"`
}

type Database interface {
	GetItems() ([]*Item, error)
	GetRecipe(id string) (*Recipe, error)

	// TOOD(lexmach): add factory/belt getters
}
