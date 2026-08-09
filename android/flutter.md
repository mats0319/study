# flutter

[flutter doc](https://docs.flutter.cn/install/quick)
[flutter实战](https://book.flutterchina.club)
[flutter codelab](https://codelabs.developers.google.com/codelabs/flutter-codelab-first?hl=zh-cn#0)

## 常用命令

- `flutter build appbundle`
- `flutter build apk --split-per-abi`
- `flutter install --use-application-binary=build/app/outputs/flutter-apk/app-arm64-v8a-release.apk`
    - 不要像flutter文档中提到的，直接使用`flutter install`，因为默认会下载`app-release.apk`，
      只有当你将所有架构的包都打在一起的时候（`build`不使用`--split-per-abi`参数），这样做才符合预期
- `adb shell getprop | grep cpu` 查看手机cpu架构：（需要手机开启usb调试）

## codelab

flutter文档中提供了一个codelab，完成之后对flutter编程有了一些基本的了解：

- 在使用过程中大概了解了哪些Widget有长宽、padding/margin等字段可以设置
- 用flutter写样式整体上比起使用html+css更容易
    - 举个例子：html style会覆盖class、css !important又会覆盖style。
      类似的还有很多，总的来说，当你写html的时候，很容易*写了居中但是没有居中*，
      flutter把样式封装的很好，基本上可以做到所见即所得

## 技术点

### 防止连续点击

如果不防止连续点击可能发生什么：

- create函数重复调用，第二次开始报错：id已存在（或者创建多个相同实例，取决于代码实现）

思路：

1. 节流：一定时间内的点击只执行一次
2. 防抖：设置延迟，延迟期间没有点击再执行
3. 使用变量控制：类似锁，事件需要拿到锁才能执行，执行之后释放锁

解决办法：执行期间将按钮变成loading

## widget的状态

有/无状态的widget：（stateful/stateless）指widget是否需要跟随指定条件而变化（例如用户输入）
（如果需要在build以外使用context，也需要使用有状态的widget）

### statelessWidget

生命周期函数只有两个：`createElement`/`build`

### statefulWidget

生命周期函数有三组：

- 初始化阶段：`createState`/`initState`
    - `createState`：抽象方法，必须重写该方法
    - `initState`：对应Android的`onCreate`方法/ios的`viewDidLoad`方法，常在其中进行一些初始化操作
- (重新)渲染阶段：`didChangeDependencies`/`build`/`didUpdateWidget`
    - `didChangeDependencies`：创建阶段调用，InheritedWidget 相关
    - `build`：页面每次渲染都会调用，使用`setState`也会调用
    - `didUpdateWidget`：只有在父容器组件重绘时调用
- 销毁阶段：`deactivate`/`dispose`
    - `dispose`：常在其中进行一些资源的释放与销毁操作

## 异常处理

可以使用dart支持的`try-catch-finally`，也可以使用`method().then().catchError().whenComplete()`
