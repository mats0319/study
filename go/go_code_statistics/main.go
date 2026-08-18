package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type statisticalResult struct {
	goFiles []*fileItem
	goTests []*fileItem
}

type fileItem struct {
	name  string
	lines int
}

func (sr *statisticalResult) traverseDir(dir string) {
	entry, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for i := range entry {
		if strings.HasPrefix(entry[i].Name(), ".") {
			continue // ignore shadow files and folders, they usually not go files
		}

		name := mustDir(dir) + entry[i].Name()
		if entry[i].IsDir() {
			sr.traverseDir(name)
			continue
		}

		switch {
		case strings.HasSuffix(entry[i].Name(), "_test.go"):
			sr.goTests = append(sr.goTests, &fileItem{
				name:  name,
				lines: statisticsFileLines(name),
			})
		case strings.HasSuffix(entry[i].Name(), ".go"):
			sr.goFiles = append(sr.goFiles, &fileItem{
				name:  name,
				lines: statisticsFileLines(name),
			})
		}
	}
}

// 不需要手动排序，os.readDir已经是排序过后的结果了
func (sr *statisticalResult) print() {
	var lineCountSummary int

	log.Println("> Go Files: ")
	lineCount := 0
	for _, file := range sr.goFiles {
		log.Printf("- %s, lines: %d\n", file.name, file.lines)
		lineCount += file.lines
	}
	log.Printf("- Summary, Files: %d, Lines: %d\n", len(sr.goFiles), lineCount)
	log.Println()

	lineCountSummary = lineCount
	lineCount = 0

	log.Println("> Go Test Files: ")
	for _, file := range sr.goTests {
		log.Printf("- %s, lines: %d\n", file.name, file.lines)
		lineCount += file.lines
	}
	log.Printf("- Summary, Files: %d, Lines: %d\n", len(sr.goTests), lineCount)
	log.Println()

	lineCountSummary += lineCount

	log.Printf("> Summary, Files: %d, Lines: %d\n", len(sr.goFiles)+len(sr.goTests), lineCountSummary)
}

func main() {
	counterIns := &statisticalResult{}
	counterIns.traverseDir("./")
	counterIns.print()
}

func statisticsFileLines(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // 一般来说go代码一行不会超过64KB ^_^
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	err = scanner.Err()
	if err != nil {
		panic(err)
	}

	return lineCount
}

func mustDir(dir string) string {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	return dir
}
