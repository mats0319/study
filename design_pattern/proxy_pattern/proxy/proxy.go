package mario

import (
	"fmt"
	"math/rand"
)

type Proxy struct {
	Employees            []*EmbroideredGirl
	Random               *rand.Rand
	CustomizableCustomer []string
	Counter              int
}

func (p *Proxy) Embroider(size string) string {
	index := p.Random.Intn(len(p.Employees))

	p.Counter++

	return p.Employees[index].Embroider(size)
}

func (p *Proxy) EmbroiderCustomized(customer, size, requirements string) string {
	var flag bool
	for i := range p.CustomizableCustomer {
		if customer == p.CustomizableCustomer[i] {
			flag = true
			break
		}
	}

	if !flag {
		return fmt.Sprintf("Customize require from %s is denied because of permission.", customer)
	}

	index := p.Random.Intn(len(p.Employees))

	assignRes := p.Employees[index].EmbroiderCustomized("", size, requirements)

	p.Counter++

	return fmt.Sprintf("Proxy received a customized from %s, res: %s", customer, assignRes)
}
