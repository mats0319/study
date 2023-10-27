package test

import (
	"fmt"
	"regexp"
	"testing"
)

func TestServiceRE(t *testing.T) {
	str := []byte(`const URI_GetList = "/list/get"`)

	re := regexp.MustCompile(`const\s+URI_(\w+)\s*=\s*"([\w/]+)"`)
	reRes := re.FindAllSubmatch(str, -1)
	// reRes[0][1]: first '()' reRes[0][2]: second '()'
	// reRes 第一层：str中所有匹配的字符串 第二层：第一个元素是字符串整体，第二个元素开始是括号提取的子串 第三层：[]byte表示字符串

	fmt.Println(string(reRes[0][1]), string(reRes[0][2]))
}

func TestMessageRE(t *testing.T) {
	// type Playlist struct {
	// 	 ID          string  `json:"id" yaml:"id"`
	// 	 Name        string  `json:"name" yaml:"name"`
	// 	 Description string  `json:"description" yaml:"description"`
	// 	 MusicList   []Music `json:"music_list" yaml:"music_list"`
	// }

	var str = []byte("type Playlist struct {\n\tID          string  `json:\"id\" yaml:\"id\"`\n\tName        string  `json:\"name\" yaml:\"name\"`\n\tDescription string  `json:\"description\" yaml:\"description\"`\n\tMusicList   []Music `json:\"music_list\" yaml:\"music_list\"`\n}")

	msgRE := regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	msgRERes := msgRE.FindAllSubmatch(str, -1)

	fmt.Println(string(msgRERes[0][1]) /*, string(msgRERes[0][2])*/)

	fieldRE := regexp.MustCompile(`\w+\s+([\[\]\w]+)\s+.*json:"(\w+)".*`)
	fieldRERes := fieldRE.FindAllSubmatch(msgRERes[0][2], -1)
	for i := range fieldRERes {
		fmt.Println(string(fieldRERes[i][1]), string(fieldRERes[i][2]))
	}
}
