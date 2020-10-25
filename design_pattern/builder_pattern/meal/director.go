package mario

type Director struct {
	Builders     []BuilderI
	builderIndex int
}

var _ MakeSetMealI = (*Director)(nil)

func (d *Director) MakeSetMeal() (res SetMeal) {
	res = SetMeal{
		Food:  Item{},
		Drink: Item{},
	}

	if d.builderIndex < 0 || d.builderIndex >= len(d.Builders) {
		return
	}

	res.Food = d.Builders[d.builderIndex].MakeFood()
	res.Drink = d.Builders[d.builderIndex].MakeDrink()

	return
}

func (d *Director) SetBuilderIndex(i int) {
	d.builderIndex = i
}
