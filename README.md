# 学习笔记

## 目录

- go
    - [x] doc
    - [x] encrypt transmission：加密传输，可以让两个人在不安全的通信信道上安全的传递信息
    - [x] generate avatar：生成像素风格对称图片，可用作默认头像
    - [x] gocts：根据go语言定义的接口结构，生成相应的ts发送http请求的代码
    - [x] utils
    - [x] listen_bilibili：使用B站作为音源的听歌工具，[代码地址](https://github.com/mats0319/listen_bilibili)
    - [x] unnamed_plan[代码地址](https://github.com/mats0319/unnamed_plan)
- vue3
    - [x] vue doc：官方文档阅读笔记
    - [x] vue ecology：vue生态，包括但不限于pinia、vite config、ts编译选项、eslint
    - [x] vue router
- android
    - [x] dart：简单过一遍dart语法，知道dart代码大概是什么样子的
    - [x] flutter：学习flutter官方文档提供的codelab
    - [x] totp：使用dart+flutter编写的符合RFC 6238标准的应用，[代码地址](https://github.com/mats0319/totp)
- 笔记
    - [x] git/github使用记录
    - [x] JWT学习笔记
    - [x] make学习笔记
    - [x] markdown学习笔记
    - [x] 正则表达式学习笔记
    - [x] shell学习笔记
    - [x] totp学习笔记

## 为release计算hash

如果一个项目的编译过程比较复杂，通常考虑将编译好的内容发布到release，而发布编译好的内容，又通常需要一个hash防伪。
记录如何计算一个文件/文件夹的hash：

- 计算文件hash：
    - linux：`sha1sum [file name]`
    - windows: `Get-FileHash -Path [file name] -Algorithm SHA1`
- 计算文件夹hash：
    - linux：`find [folder path] -type f -print0 | xargs -0 sha1sum | sha1sum`
    - windows：`Get-ChildItem -Path [folder path] -File -Recurse | ForEach-Object { Get-FileHash $_.FullName 
      -Algorithm SHA1 } | Sort-Object Path | ForEach-Object { $_.Hash } | Get-FileHash -Algorithm SHA1`
    - windows也可以装一个git，然后使用linux的命令

计算文件夹hash，它的过程是找到文件夹里的所有文件（排除文件夹、链接等），计算每个文件的hash，
然后再对这些hash计算一次hash，得到的就是文件夹的hash

## 代码库重构

因为以下原因，决定重新调整代码库的组织结构：

- .git文件夹内容过多（一个41M的仓库，.git文件夹占40M）
- 代码/工具创建在了不适合的位置或随着逐渐的开发不适合放在最初的位置了

经过了解，决定使用创建孤儿分支+gc的方式，原地删除仓库全部历史，操作步骤见笔记-github使用，考虑编写shell脚本

现有的代码、工具、其他内容以及想要创建的内容：

- lc：lc题库学习
- lb：app
- totp：app
- up：a web project
- study-note：学习笔记
- study-android/go/vue3：编程语言学习
- study-go-st：加密传输工具
- study-go-gocts：go-ts之间传递消息的结构生成工具
- 算法导论代码实现

新的个人名下代码仓库及其功能/内容简介：

- mats0319：同名仓库，github个人主页展示，包含当前study大部分功能：介绍其他代码库的基本情况、收录笔记
- study：学习，不知道放在哪的内容都可以先放在这里
- lc、lb、totp、up、st：已经归档的内容、需要提供release的内容、或者已经有明确开发方向且体量不适合和其他内容混在一起的内容

重构步骤：

- 建库，可用资源移动
- go代码调整go.mod module名称，与新的组织结构一致
- 调整文字材料：Readme.md
- android项目重新创建
