# Go + eth

## 介绍

本节主要以go语言为主，额外涉及到编译sol合约等问题

1. 使用go语言向eth链发交易
2. 一些eth链的功能测试：
    1. 将多笔交易打成一笔，更节省gas费
    2. 跨合约调用需要消耗更多的gas

todo：拟添加扫描链的代码

## 体验

当前代码即可体验，其中infura key、地址(0xe955)是我的

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
