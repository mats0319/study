package mario

type Decorator struct {
	Finery
}

func (d *Decorator) Decorate(f Finery) {
	d.Finery = f

	return
}
