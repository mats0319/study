package mario

import "fmt"

type Person struct {
	Name string
}

func (p *Person) Show() {
	if p.Name == "" {
		p.Name = "Mario" // 注意：在方法内对接收者值的修改，仅在方法内生效，既：即使执行了这一行代码，离开方法之后，p.Name还是空字符串
	}

	fmt.Printf("我叫: %s, 我穿着:", p.Name)

	return
}
