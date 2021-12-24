package mario

import "fmt"

type Client struct {
	Name   string
	Server *Server
}

var _ ClientI = (*Client)(nil)

func (c *Client) Update() {
	fmt.Println(fmt.Sprintf("%10s received notification from server, new status: `%s`", c.Name, c.Server.ServerStatus))
}

func (c *Client) GetName() string {
	return c.Name
}
