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
	mlog.Initialize(mlog.W_File)
	defer mlog.Close()
	initStdoutLogger()

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
	inputRE := regexp.MustCompile(`(\w+)`) // 匹配第一个字符串，大小写不敏感 (== [0-9A-Za-z_])
	scanner := bufio.NewScanner(os.Stdin)

ALL:
	for { // block
		info("Enter Your Command ('h' for help): ")

		if !scanner.Scan() {
			break
		}

		text := strings.ToLower(strings.TrimSpace(scanner.Text()))
		matched := inputRE.FindString(text)
		switch matched {
		case "h", "help":
			printHelp()
		case "g", "gen", "generate":
			err := internal.GenerateKeyPair(false)
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
			info(fmt.Sprintf("Unknown input: '%s', 'h' for help.\n", text))
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

var logger *slog.Logger

func initStdoutLogger() {
	h, e := mlog.NewHandler(mlog.W_Stdout)
	if e != nil {
		panic(e)
	}

	logger = slog.New(h)
}

func info(message string) {
	data := fmt.Sprintf("> %s", message)
	mlog.Log(logger, slog.LevelInfo, data)
	mlog.Info(data)
}

func printResult(message string, err error) {
	if err != nil {
		data := fmt.Sprintf("- %s Failed, %s\n", message, err.Error())
		mlog.Log(logger, slog.LevelError, data)
		mlog.Error(data)
	} else {
		data := fmt.Sprintf("> %s Success.\n", message)
		mlog.Log(logger, slog.LevelInfo, data)
		mlog.Info(data)
	}
}
