package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/mats0319/secure_transfer/internal"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

func main() {
	mlog.Initialize()
	defer mlog.Close()

	workDir()

	start()
}

func workDir() {
	path, err := os.Getwd()
	if err != nil {
		printResult("Get Current Path", err)
	}
	if !strings.HasSuffix(path, "/") { // must dir
		path += "/"
	}
	info("Work Dir: " + path)
}

func start() {
	scanner := bufio.NewScanner(os.Stdin)

ALL:
	for { // block
		info("Enter Your Command ('h' for help) .")

		if !scanner.Scan() {
			break
		}

		text := strings.ToLower(strings.TrimSpace(scanner.Text()))
		// 匹配第一个字符串，大小写不敏感 (== [0-9A-Za-z_])
		matched := regexp.MustCompile(`(\w+)`).FindString(text)
		switch matched {
		case "h", "help":
			printHelp()
		case "g", "gen", "generate":
			err := internal.GenerateKeyPair()
			printResult("Generate Key Pair", err)
		case "i", "init", "initialize":
			err := internal.InitMessageFile()
			printResult("Initialize Message File", err)
		case "e", "encrypt":
			err := internal.Encrypt()
			printResult("Encrypt", err)
		case "d", "decrypt":
			err := internal.Decrypt()
			printResult("Decrypt", err)
		case "exit", "q":
			info("Exit.")
			break ALL
		default:
			info("Unknown input: \"" + text + "\", 'h' for help.")
		}
	}
}

func printHelp() {
	info(`Options:
  - h: this help
  - g: generate public & private key into files ('./priv.key' & './PUB.KEY')
  - i: initialize message file ('./message.txt')
  - e: encrypt plain text from './message.xxx' and write cipher to './CIPHER.XXX'
  - d: decrypt cipher from './CIPHER.XXX' and write plain text to './message_decrypted.xxx'
  - exit: exit program
`)
}

var logger = mlog.DefaultLogger()

func info(message string) {
	mlog.Log(logger, slog.LevelInfo, fmt.Sprintf("> %s", message))
}

func printResult(message string, err error) {
	if err != nil {
		mlog.Log(logger, slog.LevelError, fmt.Sprintf("%s Failed, %s\n", message, err.Error()))
	} else {
		mlog.Log(logger, slog.LevelInfo, fmt.Sprintf("> %s Success.\n", message))
	}
}
