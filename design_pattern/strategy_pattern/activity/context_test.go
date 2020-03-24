package mario

import (
	"fmt"
	"log"
)

func ExampleStrategyContext_CalculateSummary() {
	sc := &StrategyContext{}

	result, err := sc.CalculateSummary(&ActivityNormal{ActivityDetails{Summary:100}})
	if err != nil {
		log.Fatalln("计算错误：", err)
	}

	fmt.Println(result)

	return

	// OutPut:
	// 100
}
