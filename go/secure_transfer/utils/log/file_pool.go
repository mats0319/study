package mlog

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

type fileRef struct {
	path string
	f    *os.File
	refs int
}

var (
	filePoolMu sync.Mutex
	filePool   = make(map[string]*fileRef) // 绝对路径 -> 共享句柄
)

func getFileRef(path string) (*fileRef, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Println("get absolute path failed, error: ", err)
		return nil, err
	}

	filePoolMu.Lock()
	defer filePoolMu.Unlock()

	ref, ok := filePool[abs]
	if ok { // 文件已被其他writer使用，复用相同文件句柄
		ref.refs++
		return ref, nil
	}

	f, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	ref = &fileRef{path: abs, f: f, refs: 1}
	filePool[abs] = ref

	return ref, nil
}

func (ref *fileRef) close() {
	filePoolMu.Lock()
	defer filePoolMu.Unlock()

	ref.refs--
	if ref.refs <= 0 {
		_ = ref.f.Close()
		delete(filePool, ref.path)
	}
}
