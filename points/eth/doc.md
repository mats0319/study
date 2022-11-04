# 以太坊相关

对BTC链了解较少，不做整理
对波卡链了解停留在老版本，不做展开
工作中主要涉及的是ETH链，所以重点整理ETH链相关内容

## 介绍

1. visualization tool：合约可视化工具
    1. 一个合约操作界面，可以部署合约、调用合约函数
    2. 详见路径下`doc.md`文档
2. test：一些eth链的功能测试
    1. 将多笔交易打包成一笔，更节省gas费
    2. 跨合约调用需要消耗更多的gas
3. contracts：放一些我写的合约
    1. 暂时只放工具库，因为项目没有上线

## require

[reference](https://goethereumbook.org/smart-contract-compile)

1. `.sol`智能合约编译工具、部署代码生成工具
    1. 也可以使用其他工具编译，如truffle，hardhat
2. `abigen`：根据合约代码生成go调用代码的工具

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
