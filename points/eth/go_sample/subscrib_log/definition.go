package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type EventParams_TestEventName struct {
	Admin       common.Address
	InvokeTimes uint8
}

func (ep *EventParams_TestEventName) String() string {
	return fmt.Sprintf("Admin: %s, Invoke times: %d", ep.Admin.Hex(), ep.InvokeTimes)
}
