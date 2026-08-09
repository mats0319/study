package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var logList *widget.List

func Bottom() *fyne.Container {
	title := widget.NewLabel("> Operate Log: ")

	return container.NewVBox(widget.NewSeparator(), title, newLogList())
}

func newLogList() *fyne.Container {
	logList = widget.NewList(
		func() int { return len(logData) },
		func() fyne.CanvasObject { return widget.NewLabel("log placeholder") },
		func(index widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(logData[index].String())
		},
	)

	background := canvas.NewRectangle(color_MainDark)
	background.SetMinSize(fyne.NewSize(800, 200))

	return container.NewStack(background, logList)
}
