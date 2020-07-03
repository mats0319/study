package mario

import "fmt"

type EmbroideredGirl struct {
    Name string
}

func (eg *EmbroideredGirl) Embroider(size string) string {
    return fmt.Sprintf("%s embroidered a product with size: %s", eg.Name, size)
}

func (eg *EmbroideredGirl) EmbroiderCustomized(_, size, requirements string) string {
    return fmt.Sprintf("%s embroidered a customized product with size: %s, requirements: %s", eg.Name, size, requirements)
}
