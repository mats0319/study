# Go module

go语言的依赖管理系统

通过`go.mod`文件定义，文件描述了module的属性，包括它的依赖和使用的go版本

## 初始化

`go mod init [module path]`

1. 要求当前路径不存在`go.mod`文件
2. 会根据import comments或版本控制配置猜测module path，也可以显式指定module path

## 常用命令

1. `go mod tidy`：整理、下载依赖
2. `go mod download`：下载依赖

## go.mod 文件结构

### module

当前module的path，供自身或其他module引用，通常与module路径相同

module path可以结合版本信息一起使用，例如`github.com/go-pg/pg/v10`和`github.com/go-pg/pg`就是不同的代码

### go

表示当前代码包要求的最低go版本

### require

指定依赖库及其最小版本，版本格式应为**release格式**或**go生成的伪版本号格式**(Pseudo-version number)

1. release格式版本号举例：`v1.2.3`
2. go生成的伪版本号举例：`v0.0.0-20170915032832-14c0d48ead0c`

### replace

通过另一个版本或本地目录，替换指定依赖的一个版本或全部版本

#### 语法

`replace <module-path> [module version] => <replacement-path> [replacement version]`

1. module version可选，不填表示替换该依赖的全部版本
2. replacement path，要求为go module管理的代码
3. replacement version则根据replacement path决定：
    1. 若replacement path是远程目录，则需要replacement version
       1. e.g. `github.com/mats9693/utils v0.0.0`
    2. 若replacement path是本地目录，则不要replacement version
       1. e.g. `../utils`

#### 使用场景

1. 应用本地代码
2. 应用fork的代码
3. 代码版本变化

### exclude

排除一个module的指定版本或全部版本，一般用于排除一个checksum无效的版本

## 版本规则

标准格式：`@<version>`，例如：`golang.org/x/text@v0.3.0`

1. 版本可以简写，例如`@v1`表示v1的最后一个可用版本
2. `commit hash`：如果代码在版本管理系统中，可以使用commit hash、分支名或其他版本管理系统中适用的唯一标识，例如：
    1. `go get golang.org/x/text@master`：master分支上的最新提交
    2. `go get github.com/go-pg/pg/v10`和`go get github.com/go-pg/pg@v10`是不一样的，因为它们指向不同的go module
3. 明确指定下载一个更低版本的依赖，可以对该依赖进行**降级**
4. 版本后缀`@none`表示删除依赖
5. 版本后缀`@latest`表示最新次要版本，即`v1.2.3`中的`2`
6. 若不提供版本信息，则使用最后一个tag、发布(release)版本或最新的提交（优先级从前往后依次降低）

# 扩展：go语言推荐的版本号格式与约定

v1.4.0-beta.2

[reference](https://go.dev/doc/modules/version-numbers)

## 解释

1. `1`：主(major)版本号
2. `4`：次要(minor)版本号
3. `0`：修补(patch)版本号
4. `beta.2`：预发布标识符(pre-release identifier)

## 约定

1. 主版本号变更，表示**向后不兼容的公共API更新**(backward-incompatible public API changes)，不保证兼容历史主版本
    1. v0版本表示代码仍在开发中，不保证其稳定性和兼容性
    2. v1及以上版本的代码是稳定的，除了预发布版本
    3. “有责任心的开发者会把v0版本代码升级到v1或以上”
    4. 从v1开始，谨慎升级主版本号，因为主版本号的升级对调用者意味着重大的破坏性更新，包括API不兼容等
    5. 主版本号升级，module path也会改变
2. 次要版本号变更，表示**向后兼容的公共API更新**(backward-compatible public API changes)，保证稳定性和对历史次要版本的兼容性
    1. 次要版本号升级，会改变公共API，但调用者通常不需要修改代码
    2. 当需要调整依赖、新增方法时，可以升级次要版本
3. 修补版本号变更，表示**对公共API更新无影响**(changes that don't affect the module's public API)，保证稳定性和对历史版本的兼容性
    1. 修补版本号升级，不会改变公共API，调用者不需要修改代码
    2. 当修改bug时，可以升级修补版本
4. 预发布标识符：表示这是一个**预发布里程碑**(pre-release milestone)，不保证稳定性

# 扩展：go生成的伪版本号(Pseudo-version number)

v0.0.0-20170915032832-14c0d48ead0c

go工具会为没有tag的module（例如github的每个commit）生成一个伪版本号，专门用于go module

## 解释

两个横线将伪版本号分成三个部分

1. 基础版本前缀(baseVersionPrefix)
2. 时间：`yyyymmddhhmmss`
3. 修订标识(revisionIdentifier)：commit hash的前12位，在子版本中，是用零填充的子版本号
