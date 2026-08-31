package mlog

import (
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

const fileName = "log.log"

type handlerWriter struct {
	writerFlag writerFlag // write to 'logfile' and/or 'Stdout'

	writer io.Writer
	mu     sync.Mutex

	// 'W_File'
	file *fileRef
}

func newHandlerWriter(wf writerFlag) (*handlerWriter, error) {
	if wf <= 0 {
		log.Println("invalid writer flag: ", wf)
		return nil, errors.New("invalid writer flag")
	}

	hw := &handlerWriter{writerFlag: wf}

	writers := make([]io.Writer, 0, 2)
	if wf&W_File > 0 {
		ref, err := getFileRef(fileName)
		if err != nil {
			log.Println("open log file failed, error:", err)
			return nil, err
		}

		hw.file = ref
		writers = append(writers, ref.f)
	}
	if wf&W_Stdout > 0 {
		writers = append(writers, os.Stdout)
	}

	switch len(writers) {
	case 1:
		hw.writer = writers[0]
	case 2:
		hw.writer = io.MultiWriter(writers...)
	}

	return hw, nil
}

func (hw *handlerWriter) Write(log []byte) (err error) {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	_, err = hw.writer.Write(log)

	return
}

func (hw *handlerWriter) close() {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	if hw.file != nil {
		hw.file.close()
		hw.file = nil
	}
}
