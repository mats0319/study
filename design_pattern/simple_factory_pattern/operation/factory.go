package mario

import "log"

type OperationFactory struct {
	OperationNumber
	Operate string
}

func (o *OperationFactory) NewOperation() Operation {
	var result Operation

	switch o.Operate { // sfp1: o.Operate
	case OperateAdd:
		result = &OperationAdd{o.OperationNumber}
	case OperateSub:
		result = &OperationSub{o.OperationNumber}
	case OperateMul:
		result = &OperationMul{o.OperationNumber}
	case OperateDiv:
		result = &OperationDiv{o.OperationNumber}
	default:
		log.Fatalln("未知的运算符类型：", o.Operate)
	}

	return result
}
