# kotlin调用rust代码示例与rust加解密基准测试

目录介绍：

- native：简单的rust demo，为了跑通kotlin调用rust代码的流程
- crypto：rust加解密lib，包含功能测试代码和基准测试代码
    - core/core_test.rs：功能测试代码，保证每次调整后能够正确加密-解密
    - core/encrypt_v1.rs：基准测试比较代码，用于比较两个版本的代码性能
- viper_crypto：最终使用的代码，在crypto最高性能代码的基础上，删除了嵌入源代码的测试代码
- android：kotlin android app，用于验证流程

## 算法

kotlin传文件路径，rust打开文件，流式加密

- 流式加密：一帧一帧的读取、加密、写入
    - 帧：默认一帧大小为1M
- 流水线模型：读取、加密、写入分别在不同的线程中进行
- 流水线模型-并发加密：多个线程同时加密（最多4线程）

## 开发计划

因为个人rust水平原因，为了保持代码足够清晰，一些地方处理的可能不够优雅，例如：

- 在线程之间发送额外的空白帧表示结束，这样做可以不用单独处理空白文件的情况，但是所有情况都要发一个空白帧

还有一些没有实现的功能，因为各种各样的原因：

- 功能：根据文件大小选择不同的加密方式/加密线程数  
  未实现原因：其实应该按照文件大小、终端性能等多因素共同作用的结果选择加密方式/加密线程数，边界条件没想好怎么写
- 功能：根据文件开头写入的信息，匹配不同的解密策略  
  未实现原因：这是为了兼容多个解密策略同时存在的情况，等需要兼容的时候再写
- 功能：每一帧都会buffer.resize ()，应尝试消除这一部分memset  
  未实现原因：对应buffer.spare.capacity_mut ()，需要使用unsafe代码，我在go都没写过unsafe...

## 使用

```cmd
rustup target add aarch64-linux-android
rustup target add x86_64-linux-android  # ide模拟器使用

// cargo install cargo-ndk

cargo ndk -t x86_64 build --release  # 得到：target/x86_64-linux-android/libnative.so
```

以上仅包含部分命令行命令，其他的像是代码写法，请自行参考实际代码

```cmd
# 一些过程中可能会使用的命令
adb push ./message.txt /sdcard/Download/  # 向ide模拟器存文件
adb pull /sdcard/Download/ .  # 将模拟器的文件夹拷贝出来
adb shell ls -l /sdcard/Download/  # 查看模拟器指定目录的文件
dd if=/dev/urandom of=test_10g.bin bs=1G count=10  # 创建测试用文件
sha1sum "xxx" | cut -d" " -f1 > "xxx.sha1"  # 通过hash算法比较经过加解密的文件和原始文件是否相同
cargo test -- --nocapture  # rust会拦截部分输出，禁用之后可以在控制台看到print的内容
perf stat cargo bench  # 展示更详细的基准测试结果
```

``` rust
use std::time::Instant;

let start = Instant::now();
// ... 加密处理 ...
println!("Encryption completed in {:?}", start.elapsed());
```

使用app测试成功：

```cmd
dev0319@dev0319-MS-7E24:~/document/code/kotlin_rust/android/app/src/main/java/com/example/demo/Download$ adb shell ls -l /sdcard/Download
total 512032
-rw-rw---- 1 u0_a200 media_rw 104857600 2026-08-04 16:50 selected_2143743037.txt
-rw-rw---- 1 u0_a200 media_rw 104859766 2026-08-04 16:51 selected_3133289981.txt
-rw-rw---- 1 u0_a200 media_rw 104857600 2026-08-04 16:51 selected_decrypted_3133289981.txt
-rw-rw---- 1 u0_a200 media_rw 104859766 2026-08-04 16:50 selected_encrypted_2143743037.txt
-rw-rw---- 1 u0_a200 media_rw 104857600 2026-07-31 10:14 test_100m.txt
dev0319@dev0319-MS-7E24:~/document/code/kotlin_rust/android/app/src/main/java/com/example/demo/Download$ adb shell sha1sum "/sdcard/Download/selected_2143743037.txt"
efa52580b6e9ede1606d4bd2db4c64ea47e50a18  /sdcard/Download/selected_2143743037.txt
dev0319@dev0319-MS-7E24:~/document/code/kotlin_rust/android/app/src/main/java/com/example/demo/Download$ adb shell sha1sum "/sdcard/Download/selected_decrypted_3133289981.txt"
efa52580b6e9ede1606d4bd2db4c64ea47e50a18  /sdcard/Download/selected_decrypted_3133289981.txt
```

## 优化记录

每个优化应包含：

- from/to 版本，AI生成代码记为v1.0
- 优化性能：使用100M文件、100次加密测试性能，比较两个版本代码之间的性能差异，超过10%的予以保留
- 优化内容：具体执行了哪些优化
- 备注：可选，用来记录一些额外信息

说明：

1. 基准测试在linux系统上执行，不代表最终android系统性能提升

### v1.4：使用buffer池管理每一帧的内存分配

- 版本：v1.3 -> v1.4
- 性能提升：10% (480ms -> 435ms)
- 详情：这个优化的重点在于减少内存分配，这在android上的提升会更明显

### v1.3：流水线模型-多线程加密

- 版本：v1.2 -> v1.3
- 性能提升：40% (800ms -> 480ms)
- 详情：因为流水线模型每帧用时变成了单个环节用时的最大值，所以使用多线程处理该环节，可以进一步提高性能
- 备注：性能提升这么大更多的是因为测试在linux系统上运行

### v1.2：流水线模型

- 版本：v1.0 -> v1.2
- 性能提升：20% (1000ms -> 800ms)
- 详情：流式加密的每一帧包含读取、加密、写入三个步骤，原本使用顺序执行，每一帧的三个步骤执行完成才会执行下一帧；
  新的流水线模型会在每一帧的一个步骤执行完成之后就执行下一帧的该步骤，所以每帧平均耗时从sum (a,b,c)变成了max (a,b,c)；
- 备注：首先保证cargo test功能测试通过，然后执行perf stat cargo bench基准测试

### v1.1：取消使用encryptorBE32

- 版本：v1.0 -> v1.1
- 性能提升：2% (1000ms -> 980ms)
- 详情：取消使用encryptorBE32，自行管理流式加密的整个过程（读取、加密、写入）
- 备注：使用1g文件测试

## 讨论记录

2026.7.28

目前rust负责：派生密钥、执行加解密、读写文件、记录加解密算法信息

- 输入参数：文件路径、密钥片段*2、file id、aad
- 输出参数：无

经过讨论，决定算法细节、架构位置均保持不变

- 加密文件开头记录了加密算法、密钥派生算法、格式化版本号等算法相关信息，设计意图为兼容多个算法时使用
- aad目前没有传参，考虑到这是算法的一部分，还是保留该参数
- AI提到Android真实路径问题：程序拿到的不是真实路径、rust无法打开，经过与jayce讨论，认为这不是一个问题
- 接上一条，因为上一条不成立，调整为kotlin读写文件的理由不充足，保持由rust读写文件
