# Android 开发

## Android开发技术选型

常见技术栈：

- kotlin + jetpack compose
- kotlin + XML
- flutter
- React Native
- Unity

主要就是在`kotlin+jetpack compose`和`flutter`中间选择，列举我的情况：

- 编写简单的Android app，或许会涉及少许android naive（例如webview）
- 基本不考虑ios平台、非移动端平台 (web、windows、linux等)

这样来看，`flutter`的优势就不大了，而`flutter`的劣势是调用android naive的能力：

- flutter调用到android native的时候还是需要写kotlin，例如我要用到的webview
    - 我没亲手写kotlin，是因为引用的库做了这部分工作
    - 我个人对这一点还是挺在意的，flutter这回是有对应的包了，下回呢？但凡遇到一个flutter没有的包，我还是得去学kotlin
        - 例如google play、AWS等第三方SDK，通常只有kotlin版本而没有dart版本
- flutter自行渲染界面，如果遇到带有界面的原生功能（例如摄像头、地图、webview等）还需要协调原生界面和自己的界面

所以现在让我写android app，我可能会选择`kotlin+jetpack compose`，因为它全栈都是kotlin语言。

## jetpack compose

[UI组件](https://developer.android.com/develop/ui/compose/components?hl=zh-cn)

这一部分简单看几个例子，知道代码结构就可以开始写了，剩下的遇到了现翻手册或者代码

- 页面框架：Scaffold（提供顶部栏、底部栏和悬浮按钮的标准结构）
- 基础布局组件：Column、Row、Box
- 基础显示组件：Text、Image、Card、Spacer
- 交互与输入控件：Button、TextField、Checkbox、Switch
