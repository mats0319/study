# 以太坊相关

## require

[reference](https://goethereumbook.org/smart-contract-compile)

1. `.sol`智能合约编译工具、部署代码生成工具
    - 也可以使用其他工具链，如truffle，hardhat
2. `abigen`：根据合约代码生成go调用代码的工具

获得合约abi文件：

使用hardhat工具：`npx hardhat compile --force`

生成文件路径：`artifacts/contracts/demo.sol/Demo.json`

在该json文件中，取abi属性对应的value（一个json数组），保存为单独的文件（假设取名为abi.json）

安装abigen：`go install github.com/ethereum/go-ethereum/cmd/abigen@latest`

使用：`abigen --abi abi.json --out demo.go --pkg contract --type Demo`

# 如果我来做区块链开发面试（ethereum）

> 2023.11 mario, for fun
> reference:
> https://docs.soliditylang.org/zh/v0.8.20
> https://docs.openzeppelin.com/contracts/5.x/upgradeable
> https://github.com/ethereum/EIPs/blob/master/EIPS/eip-150.md ：消耗gas表格
> https://docs.google.com/spreadsheets/d/15wghZr-Z6sRSMdmRmhls9dVXTOpxKy8Y64oy9MvDZEQ/edit#gid=0 ：消耗gas表格

为什么不讨论solidity语法：

- solidity语法相对简单、明确
- 因为区块gas limit、执行消耗gas等原因，智能合约通常不会很复杂
- openzeppelin提供有主流token的安全实现（包括但不限于ERC20、ERC721、ERC1155），导致一个基础的token合约甚至不需要写几行代码

## solidity

1. 小知识点：在合约函数里记录event消耗gas、跨合约调用会消耗更多的gas、将多笔交易打包成一笔更节省gas
2. 请简单介绍智能合约的fallback函数
    - 基础：fallback函数可以不显式定义，当你调用一个合约上并不存在的函数时，会触发这个函数
    - 扩展：可能存在安全性问题，一个简单的例子：C1合约F函数的功能是**向调用者转账**，  
      如果用C2合约引用C1合约来调用F函数（定义C1合约类型变量），就会触发C2合约的fallback函数；  
      当然，这个例子比较简单，因为合约上转账通常使用`address.transfer()`，这个函数固定只有2300gas，啥也不够干的
3. 为什么尽量不要在solidity中使用时间
    - 时间戳和区块hash在一定程度上受矿工影响
4. 变量的存储方式：storage/memory
    - 基础：storage类似全局变量，生命周期为永久；memory类似局部变量，生命周期为当前函数
    - 扩展：storage方式存储的变量和memory方式存储的变量之间，不能直接赋值
5. 在solidity中可以拿到的全局变量（部分）：（可以询问*solidity代码中可以获得区块的哪些信息*）
    - block.gaslimit （ uint）： 当前区块 gas 限额
    - block.number （ uint）： 当前区块号
    - msg.data （ bytes calldata）： 完整的 calldata
    - msg.sender （ address）： 消息发送者（当前调用）
    - msg.sig （ bytes4）： calldata 的前 4 字节（也就是函数标识符）
    - msg.value （ uint）： 随消息发送的 wei 的数量
    - tx.gasprice （ uint）： 随消息发送的 wei 的数量
    - tx.origin （ address）： 交易发起者（完全的调用链）
6. solidity有三种底层调用方式：`address.call`/`address.delegatecall`/`address.staticcall`，请简单介绍它们
    - 基础：正常的编程中应避免使用它们，它们是为了**与不遵守ABI的合约对接，或者为了更直接地控制编码**  
      .call会改变msg.sender（会变成合约地址），.delegatecall不会改变msg.sender（还是调用者）
    - 扩展：这些函数接收`bytes memory`输入参数，参数应包含**调用哪个函数**以及**输入参数**，举个例子：
      ```solidity
        bytes memory payload = abi.encodeWithSignature("register(string)", "MyName");
        (bool success, bytes memory returnData) = address(nameReg).call(payload);
        require(success);
      ```
7. 合约升级
    - 基础：如果合约可能需要升级，设计之初就应该拆分成**数据合约**和**控制合约**，排列组合之下，有3种升级需求  
      升级的一种实现参考opzeppelin：请求到代理合约（P），然后由P合约选择具体执行功能的合约
8. 扩展：solidity是以32字节为单位来管理参数的，这体现在不是32字节的类型可能包含**脏高位**：
    - 不占用完整 32 字节的类型可能包含“脏高位”。这在当你访问 msg.data 的时候尤为重要 —— 它带来了延展性风险：
      你既可以用原始字节 0xff000001 也可以用 0x00000001 作为参数来调用函数 f(uint8 x) 以构造交易。
      这两个参数都会被正常提供给合约，并且 x 的值看起来都像是数字 1，  
      但 msg.data 会不一样，所以如果你无论怎么使用 keccak256(msg.data)，你都会得到不同的结果。

## solidity与以太坊虚拟机

1. solidity运行在以太坊虚拟机中
2. 如果合约过大（部署bytecode超过区块gas limit），是无法部署到链上的
    - 解决办法：拆分成多个合约、库(library)；将部分功能从链上移除。
3. 智能合约的**计算**和**存储**会消耗gas

## 区块链

1. 区块链的基本结构/区块和交易的对应关系是怎样的
    - 多笔交易打包成一个区块，多个区块组成一条链
2. 一个区块/交易包含哪些关键属性
    - 区块：区块号、时间（不可靠）、区块hash、上一个区块hash等
    - 交易：交易hash、交易状态（成功或失败）、from、to、value、data、gas等
    - tx.to：Ether转账，to是收款地址；创建合约，to是0地址；调用合约，to是合约地址（包括ERC20等的转账）
    - gas：tx fee = gas price * gas used，value是*额外*随交易发送的ether数量
    - tx data：这里主要指合约调用交易，data包含**你调用的是哪个函数**以及**函数的输入参数**   
      可以通过在后面直接拼接附加信息的方式，与链下程序交互
3. 私钥-公钥-地址
    - 私钥：随机生成，长度为32字节（64位16进制数）
    - 公钥：根据私钥，经过椭圆曲线加密算法生成，长度为64字节（128位16进制数）
    - 地址：公钥hash，结果取后20字节（40位16进制数）
4. 椭圆曲线加解密：使用公钥加密、使用配套的私钥解密
5. 椭圆曲线签名验签：使用私钥加密、使用签名原文和签名结果可以解析出签名地址（实际上是公钥）

## 与链交互

我们主要使用go和ts与链交互，不同语言的交互流程是相似的

以下，我们将交互分为读和写

读：（从链上获取数据）

1. 链提供了订阅方式获取数据，我们也可以自己一个区块一个区块的扫描，为什么我们推荐使用自己扫描的方法？
    - 订阅逻辑复杂（指需要为了保证不漏不重复而做大量额外的工作）
2. 接上一条，假设我们已经拿到了一个新的区块，我们应如何获取我们关注的交易？
    - 如果你想要关注ether转账交易，可以考虑在合约上写一个转账函数，然后跟其他合约调用一起处理
    - 对于合约调用交易，如果合约记录了事件(event)，可以直接通过事件筛选；  
      如果没有事件，也可以根据to地址筛选交易，然后解析tx data获取调用的函数和输入参数

写：（向链上发送交易）

1. 参考第一小节，使用go向链上发送交易只需要提供**你想调用的函数**和**输入参数**即可，其他的go-ethereum库可以帮你完成，你只需要做一些配置
2. 你当然也可以自己拼交易、发送jsonrpc，这样可以实现更多功能，例如代发交易——签名和发送交易的不是一个人
