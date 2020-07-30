package main

import mario "github.com/mats9693/study/design_pattern/observer_pattern/push"

func main() {
	var (
		server = &mario.Server{}

		client_1 = &mario.Client{Name: "Mario", Server: server}
		client_2 = &mario.Client{Name: "Phoenix", Server: server}
	)

	server.Attach(client_1)
	server.Attach(client_2)

	server.ServerStatus = "Finish_Update"
	server.Notify()
}
