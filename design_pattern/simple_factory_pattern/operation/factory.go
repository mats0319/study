package mario

type OperationFactory struct {
	OperationNumber
	Operate string
}

func (o *OperationFactory) NewOperation() Operation {
	var result Operation

	switch o.Operate {
	case "+": result = &OperationAdd{o.OperationNumber}
	case "-": result = &OperationSub{o.OperationNumber}
	case "*": result = &OperationMul{o.OperationNumber}
	case "/": result = &OperationDiv{o.OperationNumber}
	}

	return result
}
