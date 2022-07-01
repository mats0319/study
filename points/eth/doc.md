# 以太坊链相关

对BTC链了解较少，不做整理
对波卡链了解停留在老版本，不做展开
工作中主要涉及的是ETH链，所以重点整理ETH链相关内容

## 介绍

1. go sample：使用go语言与eth链交互
    1. 向链上发交易（包括部署合约、调用合约函数）
    2. 估计一笔交易使用的gas(gas used)，同时可以检测该交易是否能够执行成功
        1. tx status != 0 表示交易执行成功
    3. 订阅合约log
    4. 调用合约pure/view函数
2. contract manage：合约管理工具
    1. 可以部署依赖、部署合约、调用合约函数（带不带pure/view都能调）
    2. 使用wallet connect连接到钱包
    3. 使用web3js调用pure/view方法
3. test：一些eth链的功能测试
    1. 将多笔交易打包成一笔，更节省gas费
    2. 跨合约调用需要消耗更多的gas
4. contracts：放一些我写的合约
    1. 暂时只放工具库，因为项目没有上线

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
