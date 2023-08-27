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
    1. 暂时只放玩具库

## require

[reference](https://goethereumbook.org/smart-contract-compile)

1. `.sol`智能合约编译工具、部署代码生成工具
    1. 也可以使用其他工具链，如truffle，hardhat
2. `abigen`：根据合约代码生成go调用代码的工具

获得合约abi文件：

使用hardhat工具：`npx hardhat compile --force`

生成文件路径：`artifacts/contracts/demo.sol/Demo.json`

在该json文件中，取abi属性对应的value（一个json数组），保存为单独的文件（假设取名为abi.json）

安装abigen：`go install github.com/ethereum/go-ethereum/cmd/abigen@latest`

使用：`abigen --abi abi.json --out demo.go --pkg contract --type Demo`

```
