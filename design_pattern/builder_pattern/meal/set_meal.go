package mario

import "fmt"

type SetMeal struct {
	Food  Item
	Drink Item
}

type MakeSetMealI interface {
    MakeSetMeal() SetMeal
}

type Item struct {
	Name    string
	Package string
}

const (
	FoodName_Rice    = "rice"
	FoodName_Noodles = "noodles"

	FoodPackage_Rice    = "plate"
	FoodPackage_Noodles = "bowl"

	DrinkName_Tea  = "tea"
	DrinkName_Milk = "milk"

	DrinkPackage_Tea  = "cup"
	DrinkPackage_Milk = "cup"
)

func (s *SetMeal) MakeFood(name, packaged string) {
    s.Food = Item{
        Name: name,
        Package: packaged,
    }
}

func (s *SetMeal) MakeDrink(name, packaged string) {
    s.Drink = Item{
        Name: name,
        Package: packaged,
    }
}

func (s *SetMeal) String() string {
    return fmt.Sprintf("Set Meal: Food %7s packaged with %7s, Drink %7s packaged with %7s",
        s.Food.Name, s.Food.Package, s.Drink.Name, s.Drink.Package)
}
