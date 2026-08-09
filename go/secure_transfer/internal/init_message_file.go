package internal

import (
	"os"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

// InitMessageFile 约定：该函数其实也应该命名为`Generatexxx`，但是为了和`生成密钥对函数的输入指令`做出区分，
// 该函数最终使用`Initxxx`命名
//   - 解释1：程序中有`生成密钥对`、`生成消息文件`两个功能，通常情况下都应该以`Generatexxx`方式命名，
//     然而程序同时支持使用输入指令执行这两个功能，而一个`g`无法同时用作两个功能的指令，所以更改命名方式。
//   - 解释2：关于为什么更改命名方式，而不是使用类似`gkp+gmf`的多字母方式。
//     生成消息文件是辅助功能，它是因为我们固定了消息文件名而提供的辅助方案，因为它而改变`全部主要功能均只需要单个字母作为指令`
//     的规则我认为并不合适，最终选择将该辅助功能重新命名为`初始化消息文件`，暗合`这只是初始化的文件，你还需要修改`之意。
func InitMessageFile() error {
	err := os.WriteFile(utils.OriginFileName+utils.DefaultExtension, []byte("[Write your message here]"), 0644)
	if err != nil {
		e := utils.ErrInitMessageFile().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}
