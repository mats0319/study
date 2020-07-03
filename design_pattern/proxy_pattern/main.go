package main

import (
    "fmt"
    "github.com/mats9693/study/design_pattern/proxy_pattern/proxy"
    "math/rand"
    "time"
)

func main() {
    var (
        girl1 = &mario.EmbroideredGirl{Name: "EmbroideredGirl 1"}
        girl2 = &mario.EmbroideredGirl{Name: "EmbroideredGirl 2"}
        girl3 = &mario.EmbroideredGirl{Name: "EmbroideredGirl 3"}
        agent = &mario.Proxy{
            Employees:            []*mario.EmbroideredGirl{girl1, girl2, girl3},
            Random:               rand.New(rand.NewSource(time.Now().Unix())),
            CustomizableCustomer: []string{"Mario"},
        }
    )

    // situation 1: embroider 5 times, which girl finish the order is randomly
    {
        fmt.Println(agent.Embroider("185"))
        fmt.Println(agent.Embroider("185"))
        fmt.Println(agent.Embroider("185"))
        fmt.Println(agent.Embroider("185"))
        fmt.Println(agent.Embroider("185"))
        fmt.Println("Agent proxy embroidered girl, customer don't know who embroidered the clothes even it already on he's hand.")
        fmt.Println("-------")
    }

    // situation 2: customized embroider order, agent will denied customers has no permission
    {
        fmt.Println(agent.EmbroiderCustomized("Mario", "185", "white"))
        fmt.Println(agent.EmbroiderCustomized("someone else", "185", "white"))
        fmt.Println("-------")
    }

    // situation 3: summary of the two situation, there will be 6 order success in total
    {
        fmt.Printf("Count of orders: %d, is currect: %t\n", agent.Counter, agent.Counter == 6)
    }

    return
}
