# mlog日志库

一个基于slog的简易自研日志库，特点：

- 可以自定义一条日志的结构
- 可以写入多个输出（例如同时写入`stdout`和`log file`）
- 可以创建多个logger，同时写不同的内容

主要用于一些小工具，想要有多路输出又不想为此引入第三方库，就可以将这个包复制到工具里。

## 使用

见`log_test.go Example`

## 实现

日志结构：`[Time] [Level] [Code Position] log message | g.k1=v1 g.k2=v2`

- Time: `2006-01-02 15:04:05.000`
- Level: `INFO+2`(`slog.record.level.string`)
- Code Position: `xxx/xx.go:10`
