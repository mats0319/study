package mario

type StandardStruct interface {
	Prepare()
	Start()
	Stop()
}

var _ StandardStruct = (*Computer)(nil)

type Computer struct {
	CPU    *CPU
	Memory *Memory
	Driver *Driver
}

func NewComputer() *Computer {
	return &Computer{
		CPU:    &CPU{},
		Memory: &Memory{},
		Driver: &Driver{},
	}
}

func (c *Computer) Prepare() {
	c.CPU.Prepare()
	c.Memory.Prepare()
	c.Driver.Prepare()
}

func (c *Computer) Start() {
	c.Memory.Start()
	c.CPU.Start()
	c.Driver.Start()
}

func (c *Computer) Stop() {
	c.Driver.Stop()
	c.CPU.Stop()
	c.Memory.Stop()
}
