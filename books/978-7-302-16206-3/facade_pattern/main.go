package main

import (
    "fmt"
    "github.com/mats9693/study/books/978-7-302-16206-3/facade_pattern/computer"
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
