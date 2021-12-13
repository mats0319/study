package main

import (
    "fmt"
    "sync"
)

type sse struct {
    clients sync.Map // client id - channel
}

func newSSE() *sse {
    return &sse{
        clients: sync.Map{},
    }
}

func (s *sse) onConnect(userID string, ch chan *eventData) {
    s.clients.Store(userID, ch)
    s.countClientsNum(fmt.Sprintf("Client '%s' connect", userID))
}

func (s *sse) onDisconnect(userID string) {
    s.clients.Delete(userID)
    s.countClientsNum(fmt.Sprintf("Client '%s' disconnect", userID))
}

func (s *sse) addNotify(e *event) {
    for _, userID := range e.to {
        chI, ok := s.clients.Load(userID)
        if !ok {
            // todo: handle error
            continue
        }

        ch, ok := chI.(chan *eventData)
        if !ok {
            // todo: handle error
            continue
        }

        ch <- e.data
    }
}

func (s *sse) countClientsNum(text string) {
    count := 0
    s.clients.Range(func(_, _ interface{}) bool {
        count++
        return true
    })

    fmt.Printf("> %s, quantity of clients: %d\n", text, count)
}
