package mario

type Builder1 struct {
    SetMeal SetMeal
}

func (b *Builder1) MakeFood() Item {
    b.SetMeal.MakeFood(FoodName_Rice, FoodPackage_Rice)

    return b.SetMeal.Food
}

func (b *Builder1) MakeDrink() Item {
    b.SetMeal.MakeDrink(DrinkName_Tea, DrinkPackage_Tea)

    return b.SetMeal.Drink
}
