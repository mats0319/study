package utils

import (
	"github.com/pkg/errors"
	"log"
)

func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalln(errors.Wrap(err, msg))
	}
}

func CheckBool(ok bool, msg string) {
	if !ok {
		log.Fatalln(msg)
	}
}
