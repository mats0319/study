package mario

import (
	"fmt"
	"log"
)

func ExampleOperationFactory_NewOperation() {
	facIns := &OperationFactory{}

	facIns.NumberA = 1
	facIns.Operate = "+"
	facIns.NumberB = 1

	result, err := facIns.NewOperation().CalculateResult()
	if err != nil {
		log.Fatalf("计算错误：%v\n", err)
	}

	fmt.Println(result)

	return

	// OutPut:
	// 2
}
