package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Content() fyne.CanvasObject {
	displayContent := container.NewStack()

	var left *container.Split
	{
		background := canvas.NewRectangle(color_MainDark)
		background.SetMinSize(fyne.NewSquareSize(200))

		var senderEntry *fyne.Container
		{
			senderContent := makeSenderContent()
			displayContent.Add(senderContent) // set default

			senderButton := widget.NewButton("Sender", func() {
				displayContent.Objects = []fyne.CanvasObject{senderContent}
				displayContent.Refresh()
			})

			senderEntry = container.NewStack(background, container.NewCenter(senderButton))
		}

		var receiverEntry *fyne.Container
		{
			receiverContent := makeReceiverContent()

			receiverButton := widget.NewButton("Receiver", func() {
				displayContent.Objects = []fyne.CanvasObject{receiverContent}
				displayContent.Refresh()
			})

			receiverEntry = container.NewStack(background, container.NewCenter(receiverButton))
		}

		left = container.NewVSplit(senderEntry, receiverEntry)
		left.SetOffset(0.5)
	}

	mainContent := container.NewHSplit(left, displayContent)
	mainContent.SetOffset(0.25)

	return mainContent
}
