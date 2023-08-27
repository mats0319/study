package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type sse struct {
	clients sync.Map // client id - channel
}

func newSSE() *sse {
	return &sse{
		clients: sync.Map{},
	}
}

func (s *sse) testSSEHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Unsupported stream", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Credentials", "true") // for cookie

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientID := getCookie(r, cookieKey)

	ch := make(chan *eventSource)
	defer close(ch)

	s.onConnect(clientID, ch)
	defer s.onDisconnect(clientID)

	exitEventGenerator := make(chan struct{})
	defer close(exitEventGenerator)
	go s.generateEvent([]string{clientID}, exitEventGenerator)

ALL:
	for {
		select {
		case <-r.Context().Done():
			break ALL // client disconnect
		case data, ok := <-ch:
			if !ok {
				break ALL // ch closed
			}

			_, _ = fmt.Fprint(w, data.format())

			flusher.Flush()
		}
	}
}

func (s *sse) generateEvent(clientIDs []string, exit chan struct{}) {
	for {
		select {
		case _, ok := <-exit:
			if !ok {
				return
			}
		case <-time.After(3 * time.Second):
			s.pushEvent(newEvent(clientIDs, "", "", fmt.Sprintf("timestamp - %d", time.Now().Unix())))
		case <-time.After(5 * time.Second):
			s.pushEvent(newEvent(clientIDs, "", "time", fmt.Sprintf("timestamp - %d", time.Now().Unix())))
		}
	}
}

func (s *sse) onConnect(userID string, ch chan *eventSource) {
	s.clients.Store(userID, ch)
	s.countClientsNum(fmt.Sprintf("Client '%s' connect", userID))
}

func (s *sse) onDisconnect(userID string) {
	s.clients.Delete(userID)
	s.countClientsNum(fmt.Sprintf("Client '%s' disconnect", userID))
}

func (s *sse) pushEvent(e *event) {
	for _, userID := range e.receivers {
		chI, ok := s.clients.Load(userID)
		if !ok {
			// todo: handle error
			continue
		}

		ch, ok := chI.(chan *eventSource)
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

func getCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}

	return cookie.Value
}
