package generate_ts

import (
	"github.com/mats9693/study/go/goc_ts/data"
	"log"
	"os"
)

func GenerateUtils(utils *data.APIUtils, outDir string) {
	if !utils.NeedObjectToFormData { // if all bool variables are false, not create 'utils' file
		return
	}

	file, err := os.OpenFile(outDir+"utils.ts", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln("open utils file failed, error: ", err)
	}
	defer func() {
		_ = file.Close()
	}()

	content := data.Copyright

	// functions
	if utils.NeedObjectToFormData {
		content = append(content, utils.ObjectToFormData...)
	}

	_, err = file.Write(content)
	if err != nil {
		log.Fatalln("write utils file failed, error: ", err)
	}
}
