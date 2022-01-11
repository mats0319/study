package utils

import (
    "fmt"
    "github.com/pkg/errors"
    "log"
)

func CheckError(err error, msg string) {
    if err != nil {
        log.Fatalln(errors.Wrap(err, msg))
    }
}

func ShowParam(msg string, param interface{}) {
    fmt.Println(msg, param)
}
