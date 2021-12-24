package mario

type ClientI interface {
	Update()

	GetName() string
}
