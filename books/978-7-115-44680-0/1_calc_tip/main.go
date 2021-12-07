package main

import "fmt"

// 根据账单金额和小费比例，计算出小费金额并加到基础账单金额上，形成总账单金额
func main() {
	var (
		orderPrice float64 = 150
		tipRate    float64 = 5
	)

	fmt.Printf("order price: %.2f, tip rate: %.2f\n", orderPrice, tipRate)

	tip := orderPrice * tipRate / 100
	summaryPrice := orderPrice + tip

	fmt.Printf("tip: %.2f, summary price: %.2f\n", tip, summaryPrice)

	return
}
