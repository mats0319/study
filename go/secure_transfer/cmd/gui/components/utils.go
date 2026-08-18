package components

import (
	"fmt"
	"os"
	"strings"
)

func printResult(message string, err error) {
	if err != nil {
		Log(fmt.Sprintf("%s Failed, %s", message, err.Error()))
	} else {
		Log(fmt.Sprintf("%s Success.", message))
	}
}

func isFileExist(fullName string, mainName string) (hasSpecFile bool, hasFuzzFile bool) {
	entry, err := os.ReadDir("./")
	if err != nil {
		Log("Read Dir Failed.")
		return
	}

	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		fileInfo, err := entry[i].Info()
		if err != nil {
			Log("Get File Info Failed.")
			continue
		}

		switch {
		case fileInfo.Name() == fullName:
			hasSpecFile = true
		case strings.HasPrefix(fileInfo.Name(), mainName):
			hasFuzzFile = true
		}

		if hasSpecFile && hasFuzzFile {
			break // all check passed, break and return
		}
	}

	return
}
