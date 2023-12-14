package databases

type InputItem struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}
type Recipe struct {
	ID                int     `json:"recipe_id"`
	Name              string  `json:"recipe_name"`
	ItemID            int     `json:"item_id"`
	FactoryName       string  `json:"factory_name"`
	ProductionFactory float64 `json:"production_factory"` // TODO(lexmach): this is actually described as P/Q, think of rework
	FactoryId         int
	BeltName          string       `json:"belt_name"`
	BeltQuantity      int          `json:"belt_quantity"`
	InputItems        []*InputItem `json:"input_items"`
}

type Item struct {
	ID      int       `json:"item_id"`
	Name    string    `json:"item_name"`
	Recipes []*Recipe `json:"recipes"`
}

type RecipeRecursive struct {
	ID                int     `json:"recipe_id"`
	Name              string  `json:"recipe_name"`
	ItemID            int     `json:"item_id"`
	FactoryName       string  `json:"factory_name"`
	ProductionFactory float64 `json:"production_factory"` // TODO(lexmach): this is actually described as P/Q, think of rework
	FactoryId         int
	BeltName          string             `json:"belt_name"`
	BeltQuantity      int                `json:"belt_quantity"`
	Children          []*RecipeRecursive `json:"children"`
}

func (recipe *Recipe) ToRecursive() *RecipeRecursive {
	if recipe == nil {
		return nil
	}
	return &RecipeRecursive{
		ID:                recipe.ID,
		Name:              recipe.Name,
		ItemID:            recipe.ItemID,
		FactoryName:       recipe.FactoryName,
		ProductionFactory: recipe.ProductionFactory,
		FactoryId:         recipe.FactoryId,
		BeltName:          recipe.BeltName,
		BeltQuantity:      recipe.BeltQuantity,
		Children:          nil,
	}
}

type RecipeID = string
type ItemID = string
type Database interface {
	GetItems() ([]*Item, error)

	// Get item with ItemId = id
	GetItem(id ItemID) (*Item, error)

	// Get recipe with RecipeId = id
	GetRecipe(id RecipeID) (*Recipe, error)

	// Get RecursiveRecipe with RecipeID = id
	// Recursive in this context means that it will recursively
	// Find RecursiveRecipes for first recipe of inputItems
	GetRecipeRecursive(id RecipeID) (*RecipeRecursive, error)
	// TODO(lexmach): add support for non-first ???recipe???

	// TOOD(lexmach): add factory/belt impl
}
