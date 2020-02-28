package mario

import "log"

type OperationFactory struct {
	OperationNumber
	Operate string
}

func (o *OperationFactory) NewOperation() Operation {
	var result Operation

	switch o.Operate { // sfp1: o.Operate
	case OperateAdd: result = &OperationAdd{o.OperationNumber}
	case OperateSub: result = &OperationSub{o.OperationNumber}
	case OperateMul: result = &OperationMul{o.OperationNumber}
	case OperateDiv: result = &OperationDiv{o.OperationNumber}
	default:
		log.Fatalln("未知的运算符类型：", o.Operate)
	}

	return result
}

// 如果工厂类初始化函数中，switch的变量（注释：sfp1）是手动输入且要求严格的，那么可以通过定义成枚举类型来解决形似问题
// （即要求输入“Mario”，但使用者错误地输入了“mario”或“MARIO”的情况）
const (
	OperateAdd = "+"
	OperateSub = "-"
	OperateMul = "*"
	OperateDiv = "/"
)
