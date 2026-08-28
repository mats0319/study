package mlog

import (
	"io"
	"log"
	"os"
	"sync"
)

type handlerWriter struct {
	writeFlag writeFlag // write to logfile and/or Stdout

	writer io.Writer
	mu     sync.Mutex

	// 'W_File'
	fileIns *os.File
}

func newHandlerWriter(wf writeFlag) (hw *handlerWriter) {
	if wf <= 0 {
		log.Println("invalid write flag")
		os.Exit(1)
	}

	hw = &handlerWriter{writeFlag: wf}

	writers := make([]io.Writer, 0, 2)
	if hw.writeFlag&W_File > 0 {
		// 检查：此处的文件句柄应在程序退出前关闭（调用mlog.Close()），参考测试代码
		fileIns, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Println("open log file failed, error:", err)
			os.Exit(1)
		}

		hw.fileIns = fileIns

		writers = append(writers, fileIns)
	}
	if hw.writeFlag&W_Stdout > 0 {
		writers = append(writers, os.Stdout)
	}

	if len(writers) == 1 {
		hw.writer = writers[0]
	} else { // > 1
		hw.writer = io.MultiWriter(writers...)
	}

	return
}

func (hw *handlerWriter) Write(logBytes []byte) (err error) {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	_, err = hw.writer.Write(logBytes)
	if err != nil {
		return
	}

	return
}

func (hw *handlerWriter) close() {
	if hw != nil {
		_ = hw.fileIns.Close()
	}
}
