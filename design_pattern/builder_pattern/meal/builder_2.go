package mario

type Builder2 struct {
    SetMeal SetMeal
}

func (b *Builder2) MakeFood() Item {
    b.SetMeal.MakeFood(FoodName_Noodles, FoodPackage_Noodles)

    return b.SetMeal.Food
}

func (b *Builder2) MakeDrink() Item {
    b.SetMeal.MakeDrink(DrinkName_Milk, DrinkPackage_Milk)

    return b.SetMeal.Drink
}
