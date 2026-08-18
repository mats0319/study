package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/mats0319/secure_transfer/cmd/gui/components"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

func main() {
	mlog.Initialize(mlog.W_File)
	defer mlog.Close()

	a := app.NewWithID("Secure Transfer") // 不链式调用下去，因为a可以创建多个window
	a.Settings().SetTheme(&components.Theme{Theme: theme.DefaultTheme()})

	window := a.NewWindow("Secure Transfer")
	window.SetContent(makeUI(a))
	window.SetMaster() // 主窗口，关闭它即视为关闭程序；窗口大小由内部组件撑开、不主动设置（实际约为800*700）
	window.Show()

	a.Run()
}

func makeUI(a fyne.App) fyne.CanvasObject {
	return container.NewVBox(components.TopBar(a), components.Content(), components.Bottom())
}
