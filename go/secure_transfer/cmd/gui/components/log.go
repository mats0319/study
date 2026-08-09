package components

import (
	"fmt"
	"time"

	mlog "github.com/mats0319/secure_transfer/utils/log"
)

type OperateLog struct {
	Time    string
	Details string
}

var logData = make([]*OperateLog, 0)

func Log(details string) {
	newLog := &OperateLog{
		Time:    time.Now().Format("2006-01-02 15:04:05.000"),
		Details: details,
	}
	mlog.Info(newLog.String())

	logData = append(logData, newLog)

	if len(logData) > 1000 {
		logData = logData[len(logData)-1000:] // 保留最近1000条记录
	}

	logList.Refresh()
	logList.ScrollToBottom()
}

func (l *OperateLog) String() string {
	return fmt.Sprintf("> %s, %s", l.Time, l.Details)
}
