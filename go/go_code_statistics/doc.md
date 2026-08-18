# go代码行数统计工具

本工具可以统计当前目录下的所有go代码行数，并根据go代码、go测试代码进行区分。

## 使用

```cmd
go install
go_code_statistics
```

目前仅支持通过切换目录执行命令的方式改变检查路径，以下是一个运行结果示例：

```cmd
2026/08/18 15:30:49 > Go Files: 
2026/08/18 15:30:49 - ./cmd/cli/main.go, lines: 112
2026/08/18 15:30:49 - ./cmd/gui/components/bottom.go, lines: 31
2026/08/18 15:30:49 - ./cmd/gui/components/content.go, lines: 51
2026/08/18 15:30:49 - ./cmd/gui/components/data.go, lines: 63
2026/08/18 15:30:49 - ./cmd/gui/components/log.go, lines: 36
2026/08/18 15:30:49 - ./cmd/gui/components/receiver.go, lines: 63
2026/08/18 15:30:49 - ./cmd/gui/components/sender.go, lines: 74
2026/08/18 15:30:49 - ./cmd/gui/components/theme.go, lines: 27
2026/08/18 15:30:49 - ./cmd/gui/components/top.go, lines: 63
2026/08/18 15:30:49 - ./cmd/gui/components/utils.go, lines: 48
2026/08/18 15:30:49 - ./cmd/gui/main.go, lines: 29
2026/08/18 15:30:49 - ./internal/decrypt.go, lines: 97
2026/08/18 15:30:49 - ./internal/encrypt.go, lines: 90
2026/08/18 15:30:49 - ./internal/generate_key_pair.go, lines: 92
2026/08/18 15:30:49 - ./internal/init_message_file.go, lines: 26
2026/08/18 15:30:49 - ./internal/lib/decryptor_once.go, lines: 142
2026/08/18 15:30:49 - ./internal/lib/decryptor_stream.go, lines: 237
2026/08/18 15:30:49 - ./internal/lib/define.go, lines: 11
2026/08/18 15:30:49 - ./internal/lib/encryptor_once.go, lines: 136
2026/08/18 15:30:49 - ./internal/lib/encryptor_stream.go, lines: 260
2026/08/18 15:30:49 - ./internal/lib/file_header.go, lines: 202
2026/08/18 15:30:49 - ./internal/lib/stream_frame.go, lines: 105
2026/08/18 15:30:49 - ./test/utils.go, lines: 75
2026/08/18 15:30:49 - ./utils/const.go, lines: 37
2026/08/18 15:30:49 - ./utils/error.go, lines: 79
2026/08/18 15:30:49 - ./utils/errors.go, lines: 61
2026/08/18 15:30:49 - ./utils/log/handler.go, lines: 150
2026/08/18 15:30:49 - ./utils/log/handler_writer.go, lines: 136
2026/08/18 15:30:49 - ./utils/log/main.go, lines: 63
2026/08/18 15:30:49 - ./utils/utils.go, lines: 90
2026/08/18 15:30:49 - Summary, Files: 30, Lines: 2686
2026/08/18 15:30:49 
2026/08/18 15:30:49 > Go Test Files: 
2026/08/18 15:30:49 - ./test/encdec_once_test.go, lines: 40
2026/08/18 15:30:49 - ./test/encdec_stream_test.go, lines: 40
2026/08/18 15:30:49 - ./test/encrypt_bench_test.go, lines: 103
2026/08/18 15:30:49 - ./utils/errors_test.go, lines: 49
2026/08/18 15:30:49 - ./utils/log/log_test.go, lines: 60
2026/08/18 15:30:49 - Summary, Files: 5, Lines: 292
2026/08/18 15:30:49 
2026/08/18 15:30:49 > Summary, Files: 35, Lines: 2978
```
