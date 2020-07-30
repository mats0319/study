package mario

import "fmt"

var (
    _ StandardStruct = (*CPU)(nil)
    _ StandardStruct = (*Memory)(nil)
    _ StandardStruct = (*Driver)(nil)
)

type CPU struct {
}

func (c *CPU) Prepare() {
    fmt.Println("安装CPU")
}

func (c *CPU) Start() {
    fmt.Println("主板为CPU通电")
}

func (c *CPU) Stop() {
    fmt.Println("停止CPU供电")
}

type Memory struct {
}

func (m *Memory) Prepare() {
    fmt.Println("安装内存")
}

func (m *Memory) Start() {
    fmt.Println("主板为内存通电")
}

func (m *Memory) Stop() {
    fmt.Println("停止内存供电")
}

type Driver struct {
}

func (d *Driver) Prepare() {
    fmt.Println("安装硬盘")
}

func (d *Driver) Start() {
    fmt.Println("为硬盘通电")
}

func (d *Driver) Stop() {
    fmt.Println("停止硬盘供电")
}

