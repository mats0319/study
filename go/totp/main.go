package main

import (
	"log"
	"time"

	"github.com/xlzd/gotp"
)

func main() {
	otp := gotp.NewDefaultTOTP("<< Your github 2FA secret >>")

	log.Printf("> TOTP: %s, Remaining validity time: %d\n", otp.Now(), 30 - (time.Now().Unix())%30)

	time.Sleep(time.Second*30)
}
