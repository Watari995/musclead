package fooddomain

// FoodProductID は食品マスタの識別子。
type FoodProductID struct {
	// TODO: implement
}

// RegisterSource は食品データの登録元。
type RegisterSource string

const (
	RegisterSourceOpenFoodFacts RegisterSource = "open_food_facts"
	RegisterSourceUser          RegisterSource = "user"
)

// FoodProduct は食品マスタのエンティティ。
type FoodProduct struct {
	// TODO: implement
}
