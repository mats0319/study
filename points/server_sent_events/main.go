package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	listenAddr   = "127.0.0.1:9693"
	listenOrigin = "https://" + listenAddr

	cookieValidPeriod = 86400
	jwtTokenName      = "jwtToken"
	jwtSecret         = []byte("mario")
)

func main() {
	{
		http.HandleFunc("/", bindHTMLFile)
		http.HandleFunc("/cookie", testCookieHandler)
		fmt.Println("> Open web page failed, error:", exec.Command("cmd", "/c start "+listenOrigin).Start())
	}

	http.HandleFunc("/test/sse", newSSE().testSSEHandler)
	fmt.Println("> Listening at: " + listenOrigin)
	fmt.Println(http.ListenAndServeTLS(listenAddr, "server.crt", "server.key", nil))
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

	// verify jwt token and parse 'client id' from it
	clientID := ""
	{
		jwtToken := getCookie(r, jwtTokenName)
		claimsIns, err := parseToken(jwtToken)
		if err != nil {
			http.Error(w, "Parse JWT token failed", http.StatusInternalServerError)
			return
		}

		clientID = claimsIns.ClientID
	}

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

func getCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func bindHTMLFile(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	http.ServeFile(w, r, dir+"/html/index.html")
}

func testCookieHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	params := strings.Split(r.RequestURI, "?id=")
	if len(params) < 2 {
		http.Error(w, "Parse ID from URI failed", http.StatusInternalServerError)
		return
	}

	jwtToken, err := generateToken(params[1])
	if err != nil {
		http.Error(w, "Generate JWT token failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   jwtTokenName,
		Value:  jwtToken,
		MaxAge: cookieValidPeriod,
	})

	return
}
