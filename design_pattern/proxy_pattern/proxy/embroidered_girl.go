package proxy

import "fmt"

type EmbroideredGirl struct {
    Name string
}

func (eg *EmbroideredGirl) Embroider(size string) {
    fmt.Println(fmt.Sprintf("%s embroidered a product with size: %s\n", eg.Name, size))
}

func (eg *EmbroideredGirl) EmbroiderCustomized(size, requirements string) {
    fmt.Println(fmt.Sprintf("%s embroidered a customized product with size: %s, requirements: %s\n",
        eg.Name, size, requirements))
}
