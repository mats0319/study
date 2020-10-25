package main

import (
	"fmt"
	mario "github.com/mats9693/study/design_pattern/facade_pattern/computer"
)

func main() {
	facade()

	fmt.Println("-------")

	noFacade()
}

func facade() {
	var facadeIns = mario.NewComputer()

	facadeIns.Prepare()

	facadeIns.Start()

	facadeIns.Stop()
}

func noFacade() {
	var (
		CPU    = &mario.CPU{}
		Memory = &mario.Memory{}
		Driver = &mario.Driver{}
	)

	CPU.Prepare()
	Memory.Prepare()
	Driver.Prepare()

	Memory.Start()
	CPU.Start()
	Driver.Start()

	Driver.Stop()
	CPU.Stop()
	Memory.Stop()
}
