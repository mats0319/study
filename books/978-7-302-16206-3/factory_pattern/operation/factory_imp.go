package mario

type AddFactory struct {
	OperationNumber
}

func (f *AddFactory) CreateOperation() Operation {
	return &OperationAdd{f.OperationNumber}
}

type SubFactory struct {
	OperationNumber
}

func (f *SubFactory) CreateOperation() Operation {
	return &OperationSub{f.OperationNumber}
}

type MulFactory struct {
	OperationNumber
}

func (f *MulFactory) CreateOperation() Operation {
	return &OperationMul{f.OperationNumber}
}

type DivFactory struct {
	OperationNumber
}

func (f *DivFactory) CreateOperation() Operation {
	return &OperationDiv{f.OperationNumber}
}
