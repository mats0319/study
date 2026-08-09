# Go Weekly阅读笔记

记录自己在go weekly了解到的优秀文章、资源以及工具等内容。[subscribe](https://golangweekly.com/)

## 阅读

- https://beej.us/guide/bglcs/html/split/ ：计算机科学学习方法论，个人编写的学习计算机科学的宏观方法论
- https://atlas9.dev/blog/soft-delete.html ：pg软删除的缺点及替代方案（触发器或预写入日志）
- https://github.com/gibbok/typescript-book ：一本简洁的ts手册
- https://eblog.fly.dev/ginbad.html ：“gin是一个糟糕的软件库”
- https://www.sophielwang.com/blog/jpeg ：详细介绍jpeg图片格式压缩的原理
- https://medium.com/@m0t9_/lets-add-a-conditional-expression-to-go-language-3eec7783388e ： 通过为go语言添加`?:`
  条件表达式，介绍go语言编译器（词法分析、类型检查等环节）的工作过程和编辑方法
- https://github.blog/engineering/infrastructure/how-github-uses-ebpf-to-improve-deployment-safety/ ：
  github的源代码也托管在github上，这篇文章介绍了github如何使用ebpf-go解决循环依赖的问题
  （举个例子，github无法访问时，github的开发者将无法下载和上传github的源代码）
- https://blog.iamvedant.in/containers-are-not-magic-namespaces-from-scratch ：从零开始创建一个容器，介绍容器技术原理
- https://anirudhology.com/blog/building-raftly-reproducing-production-failures ： 实现一个可控的raft共识，通过删除一个环节来模拟和理解对应环节的作用
- https://corrode.dev/learn/migration-guides/go-to-rust/ ：一篇go和rust的对比文章
- https://www.alexedwards.net/blog/go-experiments-explained ：go实验性功能，它们可能默认开启或关闭、正在推进加入正式版本或已搁置
- https://segflow.github.io/post/fast-file-search-go/ ：流处理的优化过程，将基础代码（加载到内存然后按byte读取）的效率一路提高到66倍
- https://zackoverflow.dev/writing/why-does-tsgo-use-so-much-memory ：介绍ts的新编译器tsgo为什么占用这么多内存
- https://internals-for-interns.com/posts/the-go-lexer/ ：理解go编译器（系列文章）
- https://internals-for-interns.com/posts/understanding-go-runtime/ ：理解go运行时（系列文章）
- https://spf13.com/p/go-the-agentic-language/ ：从ts7编译器用go改写开始讲起，主要介绍了AI写代码go/rust代码的对比
- https://postgresisenough.dev/ ：一个介绍如何使用pg实现 *通常通过引入其他系统实现的*功能的网站， 例如全文检索
  (elasticsearch)、缓存 (redis)、文档存储 (mongoDB)... pg有这些内容的处理方案，但问题在方案之外：
  如果一个功能相对边缘化，则可以尝试使用pg；如果一个功能是核心或者有性能要求，则应考虑专门的系统
- https://opentelemetry.io/blog/2026/go-compile-time-instrumentation-v1/ ：go语言自动插桩工具，
  可以在编译期间自动注入监测代码，增加系统可观测性（监控服务响应时间、数据库查询耗时、错误堆栈、资源占用等）
- https://github.blog/security/6-security-settings-every-github-maintainer-should-enable-this-week/ ：
  github建议你开启的6项仓库安全设置
- https://github.com/EvanLi/Github-Ranking/blob/master/Top100/Go.md，github ：go相关、star最多的100个库，每天更新
- https://blog.jetbrains.com/go/2026/07/20/escape-analysis/ ，介绍go语言逃逸分析，以及goland 2026.2如何帮助进行逃逸分析
- https://rednafi.com/go/supervised-fire-and-forget/ ，一种通过goroutine池来管理goroutine的方式，
  避免其数量不受控制、无法强制退出、panic导致主进程崩溃、因主进程退出而中断等问题

## 代码/工具

- github.com/aperturerobotics/goscript：将go代码在ast级转换为ts
- github.com/hajimehoshi/ebiten：go语言2D引擎
- github.com/templui/templui：一个让人可以使用go语言编写网页的库，内置支持tailwindCss
- github.com/markel1974/godoom：go语言3d引擎
- github.com/six-ddc/plow: http基准测试工具
- github.com/gookit/validate: 通用验证器和过滤器
- github.com/fyne-io/docs.fyne.io: go GUI
