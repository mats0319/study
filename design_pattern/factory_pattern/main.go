package main

import (
    mario "github.com/mats9693/study/design_pattern/factory_pattern/operation"
    "log"
)

func main() {
    defaultONIns := mario.OperationNumber{
        NumberA: 10,
        NumberB: 5,
    }

    var ofIns mario.OperationFactory = &mario.AddFactory{OperationNumber: defaultONIns}

    res, err := ofIns.CreateOperation().CalculateResult()
    if err != nil {
        log.Fatalln("call failed, error: ", err)
    }

    log.Printf("10+5 = %f, reach expect? %t\n", res, res == defaultONIns.NumberA + defaultONIns.NumberB)

    return
}
