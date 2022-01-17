# Go + eth

## 功能介绍

1. 使用go语言与eth链交互
    1. 向链上发交易（包括部署合约、调用合约函数）
    2. 订阅合约log
    3. 调用合约view、pure函数（无gas消耗，不需要私钥）
2. 一些eth链的功能测试：
    1. 将多笔交易打包成一笔，更节省gas费
    2. 跨合约调用需要消耗更多的gas

tips：

1. 该部分仅包含go代码，保证能够正常调用，若想运行，需自行编写和部署智能合约并相应调整go代码

## 测试地址

提供两个我的地址，有一点ropsten测试链的ETH，无主网交易历史

1. 0x0f770Af04b67bB0A8e7ba9A9f9273E85553E4C4C
    1. 私钥：094be1db3580f2878cf7f9bf862cfd5c564c2776f518905c1c40445611ae3e40
    2. 助记词：until swarm manage erosion kidney clutch later blossom planet enforce forum bicycle
2. 0xD095De8Fb4E2da680D113284bcC4E0222bbAa6a2
    1. 私钥：1e742e12b1602fb8278890ebf567132ef338c0213a1d34812a29e028b3ac72b2
    2. 助记词：link awful unlock tilt trip meadow attract theme vendor doctor daughter mirror

## require

[reference](https://goethereumbook.org/smart-contract-compile)

1. ```.sol```智能合约编译工具、部署代码生成工具
    1. 也可以使用其他工具编译，如truffle
2. ```abigen```：根据合约代码生成go调用代码的工具

命令：

```cmd 
go get -u github.com/ethereum/go-ethereum
cd $GOPATH/src/github.com/ethereum/go-ethereum/
make
make devtools

solc --abi Store.sol
solc --bin Store.sol
abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=Store.go
```
