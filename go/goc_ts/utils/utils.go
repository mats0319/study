package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// WriteFile write 'content' into 'file'
func WriteFile(filename string, content []byte) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatalln("open file failed, error: ", err)
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.Write(content)
	if err != nil {
		log.Fatalln("write file failed, error: ", err)
	}

	log.Println("Generated file: ", filename)
}

// MustGoFileName make sure 'filename' is end with '.go', e.g. "demo" => "demo.go"
func MustGoFileName(fileName string) string {
	if !strings.HasSuffix(fileName, ".go") {
		fileName += ".go"
	}

	return fileName
}

// FormatStrSliceInLine e.g. ["a","b","c"] => "a, b, c"
func FormatStrSliceInLine(data []string) string {
	if len(data) < 1 {
		return ""
	}

	var resBytes []byte
	for i := range data {
		resBytes = append(resBytes, ", "+data[i]...)
	}

	return string(resBytes[2:])
}

// MustSmall make sure first char of 'str' is small-case, e.g. "MustSmall" => "mustSmall"
func MustSmall(str string) string {
	if len(str) > 0 && 'A' < str[0] && str[0] < 'Z' {
		str = string(str[0]-'A'+'a') + str[1:]
	}

	return str
}

// MustBig make sure first char of 'str' is big-case, e.g. "mustBig" => "MustBig"
func MustBig(str string) string {
	if len(str) > 0 && 'a' < str[0] && str[0] < 'z' {
		str = string(str[0]-'a'+'A') + str[1:]
	}

	return str
}

// EmptyDir del and re-make dir
func EmptyDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatalln(fmt.Sprintf("rm %s failed, error: ", dir), err)
	}

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("'mkdir' on %s failed, error: ", dir), err)
	}
}

func MustExistDir(dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("'mkdir' on %s failed, error: ", dir), err)
	}
}

// MustDir make sure 'path' is end with '/'
func MustDir(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path
}
