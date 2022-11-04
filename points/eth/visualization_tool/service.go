package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	port, err := getFreePort()
	if err != nil {
		fmt.Println("获取空闲端口失败：", err)
		return
	}

	listenAddr := fmt.Sprintf("127.0.0.1:%d", port)

	http.HandleFunc("/", bindHTMLFile)
	// auto open webpage in Windows OS
	err = exec.Command("cmd", "/c start http://"+listenAddr).Start()
	if err != nil {
		fmt.Println("自动打开网页出现错误：", err)
		fmt.Println("请手动访问：http://" + listenAddr)
	}

	fmt.Println("> Listening at: ", listenAddr)
	fmt.Println(http.ListenAndServe(listenAddr, nil))
}

func bindHTMLFile(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	path := dir + "/ui/dist" + r.RequestURI
	if r.RequestURI == "/" {
		path += "index.html"
	}
	http.ServeFile(w, r, path)
}

func getFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println("tcp listen failed", err)
		return -1, err
	}

	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil // type assert is ok
}
