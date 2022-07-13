package main

import (
	"fmt"
	"time"
)

const daySeconds = 60 * 60 * 24

func main() {
	timestamp := time.Now().Unix()

	today0 := timestamp / daySeconds * daySeconds

	doWork := today0 + 17*3600

	if doWork < timestamp {
		doWork += daySeconds
	}

	fmt.Println("current : ", timestamp, time.Unix(timestamp, 0).In(time.UTC).String(), time.Unix(timestamp, 0).String())
	fmt.Println("today 0 : ", today0, time.Unix(today0, 0).In(time.UTC).String(), time.Unix(today0, 0).String())
	fmt.Println("do work : ", doWork, time.Unix(doWork, 0).In(time.UTC).String(), time.Unix(doWork, 0).String())
}
