package test

import (
	"fmt"
	"testing"

	"github.com/mats9693/study/go/goc_ts/parse"
)

func TestServiceRE(t *testing.T) {
	str := []byte(`const URI_GetList = "/list/get"`)

	re := parse.RequestRE
	reRes := re.FindAllSubmatch(str, -1)
	// reRes[0][1]: first '()' reRes[0][2]: second '()'
	// reRes 第一层：str中所有匹配的字符串 第二层：第一个元素是字符串整体，第二个元素开始是括号提取的子串 第三层：[]byte表示字符串

	fmt.Println(string(reRes[0][1]), string(reRes[0][2]))
}

func TestMessageRE(t *testing.T) {
	// type Playlist struct {
	// 	 ID          string  `json:"id" yaml:"id"`
	// 	 FileName        string  `json:"name" yaml:"name"`
	// 	 Description string  `json:"description" yaml:"description"`
	// 	 MusicList   []Music `json:"music_list" yaml:"music_list"`
	// }

	var str = []byte("type Playlist struct {\n\tID          string  `json:\"id\" yaml:\"id\"`\n\tName        string  `json:\"name\" yaml:\"name\"`\n\tDescription string  `json:\"description\" yaml:\"description\"`\n\tMusicList   []Music `json:\"music_list\" yaml:\"music_list\"`\n}")

	msgRE := parse.StructureRE
	msgRERes := msgRE.FindAllSubmatch(str, -1)

	fmt.Println(string(msgRERes[0][1]) /*, string(msgRERes[0][2])*/)

	fieldRE := parse.StructureFieldRE
	fieldRERes := fieldRE.FindAllSubmatch(msgRERes[0][2], -1)
	for i := range fieldRERes {
		fmt.Println(string(fieldRERes[i][1]), string(fieldRERes[i][2]))
	}
}
