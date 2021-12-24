package mario

type Operation interface {
	CalculateResult() (float64, error)
}

type OperationNumber struct {
	NumberA float64
	NumberB float64
}

// 如果工厂类初始化函数中，switch的变量（注释：sfp1）是手动输入且要求严格的，那么可以通过定义成常量或变量来解决形似问题
// （即要求输入“Mario”，但使用者错误地输入了“mario”或“MARIO”的情况）
//
// 简单来说，就是尽量不要在业务代码中，出现写死的值，哪怕这个值只用在这里，也使用一个常量把它提出来、和其他的值放在一起统一管理。
const (
	OperateAdd = "+"
	OperateSub = "-"
	OperateMul = "*"
	OperateDiv = "/"
)
