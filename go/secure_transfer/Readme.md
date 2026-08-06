# 安全的信息传输工具 Secure Message Transfer Tool

本工具可以让两个人在不安全的通信信道上，以文件为载体，安全的传递信息

## 使用

> 以命令行程序使用，参考：`/cmd/cli/script/manual.md`
> 以图形化界面程序使用，参考：`/cmd/gui/components/utils.go helpText`
>
> 我们会提供一些平台的可执行程序，但我们建议你阅读源代码后自行编译。

程序会在当前目录创建和读取文件，请在一个单独的路径中启动程序。

编译：以命令行形式使用无需外部依赖 (go install)，以GUI形式使用需要安装依赖

- 我们提供的cli编译脚本写的不好、需要安装gui相关依赖，仅供参考
- fyne存在交叉编译问题：仅根据文档中要求下载的内容，无法在linux上编译windows的可执行程序，所以gui的编译脚本也仅供参考

[下载go、C编译器和系统显卡驱动程序](https://docs.fyne.io/started/quick/)

```cmd
go get fyne.io/fyne/v2@latest
go install fyne.io/tools/cmd/fyne@latest
```

### GUI

fyne依赖cgo，如果编译环境和你的系统的glibc版本不兼容，则可执行程序无法使用。  
因此我们放弃提供gui的二进制文件，仅在本小节附使用截图。

![linux](doc/linux_use.png)

## 实现

方案：使用椭圆曲线协商算法得到共享密钥，共享密钥派生得到最终加密密钥，使用对称加密算法加密

- 曲线与密钥协商算法：Curve25519 X25519
- 最终加密密钥派生算法：sha256
- 对称加密算法：aes-gcm

过程：

1. 甲生成一个临时密钥对，使用临时私钥和乙的公钥一起协商出共享密钥，共享密钥派生得到最终加密密钥
2. 甲使用最终加密密钥加密想要发送的消息，将临时公钥、nonce、密文拼在一起，发送给乙
3. 乙使用自己的私钥和临时公钥一起协商出共享密钥，派生得到最终加密密钥，结合nonce，解密密文

加密文件编码规则：

```txt
header len + 'ST' + enc method + 
salt len + base nonce len + pubK len + aad len +
salt + base nonce + pubK + aad + 
ciphertext
```

- `header len`：整个文件头（一直到正文之前）的长度
- `ST`：开头标记
- `enc method`：一次性读取到内存中加密/流式加密
- `salt`：用于密钥派生的salt
    - 本工具中长度固定为16
- `base nonce`：aes-gcm base nonce
    - 本工具中长度固定为7
    - 一次性加密nonce：base nonce后面补0凑够位数
    - 流式加密nonce：7+1+4，base nonce + 末帧标识 + 帧索引
- `pubK`：用于解密-协商密钥的公钥
    - 本工具中长度固定为32
- `aad`：加密时的关联信息
- `ciphertext`：密文
    - 流式加密密文：每一帧包含帧nonce（完整nonce的'1+4'部分）和密文

## 技术准备

[见tech_prepare.md](doc/tech_prepare.md)

## 问题

- 编写技术验证demo遇到问题：最新版本的go（1.25）和ethereum（1.16）无法匹配。  
  因为我们没有使用以太坊的曲线 (secp256k1)。放弃go-ethereum，自行编写相关过程（密钥协商、派生和加/解密）
