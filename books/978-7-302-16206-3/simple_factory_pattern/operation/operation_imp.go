package mario

import "errors"

type OperationAdd struct {
	OperationNumber
}

func (o *OperationAdd) CalculateResult() (float64, error) {
	return o.NumberA + o.NumberB, nil
}

type OperationSub struct {
	OperationNumber
}

func (o *OperationSub) CalculateResult() (float64, error) {
	return o.NumberA - o.NumberB, nil
}

type OperationMul struct {
	OperationNumber
}

func (o *OperationMul) CalculateResult() (float64, error) {
	return o.NumberA * o.NumberB, nil
}

type OperationDiv struct {
	OperationNumber
}

func (o *OperationDiv) CalculateResult() (float64, error) {
	if o.NumberB == 0 {
		return 0, errors.New("除数不能为0")
	}

	return o.NumberA / o.NumberB, nil
}
