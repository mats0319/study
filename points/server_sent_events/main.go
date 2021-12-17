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
	// for test
	{
		http.HandleFunc("/", bindHTMLFile)
		http.HandleFunc("/cookie", getCookieHandler)
		fmt.Println("> Auto start web page error:", exec.Command("cmd", "/c start "+listenOrigin).Start())
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
	// for sse
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

	// register client
	ch := make(chan *eventData)

	s.onConnect(clientID, ch)
	defer func() {
		s.onDisconnect(clientID)
	}()

	notify := r.Context().Done()
	go func() {
		<-notify
		s.onDisconnect(clientID)
	}()

	// Event generator
	{
		go func() {
			for {
				eventName := ""
				timestamp := time.Now().Unix()
				if timestamp%2 == 0 {
					eventName = "time"
				}

				s.addNotify(newEvent([]string{clientID}, "", eventName, fmt.Sprintf("timestamp - %d", timestamp)))

				time.Sleep(time.Second * 3)
			}
		}()
	}

	for {
		data := <-ch
		_, _ = fmt.Fprint(w, data.format())

		flusher.Flush()
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

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
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
