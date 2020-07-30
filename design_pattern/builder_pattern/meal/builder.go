package mario

type BuilderI interface {
    MakeFood() Item
    MakeDrink() Item
}
