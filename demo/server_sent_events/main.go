package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	listenAddr   = "127.0.0.1:9693"
	listenOrigin = "https://" + listenAddr

	cookieValidPeriod = 86400
	cookieKey         = "clientID"
)

func main() {
	{
		http.HandleFunc("/", bindHTMLFile)
		http.HandleFunc("/cookie", testCookieHandler)
		fmt.Println("> Open web page failed, error:", exec.Command("cmd", "/c start "+listenOrigin).Start())
	}

	http.HandleFunc("/test/sse", newSSE().testSSEHandler)
	fmt.Println("> Listening at: " + listenOrigin)

	err := http.ListenAndServeTLS(listenAddr, "server.crt", "server.key", nil)
	if err != nil {
		fmt.Println("listen ans serve error:", err)
	}
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

	http.SetCookie(w, &http.Cookie{
		Name:   cookieKey,
		Value:  params[1],
		MaxAge: cookieValidPeriod,
	})

	return
}
