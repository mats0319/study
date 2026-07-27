package components

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mats0319/secure_transfer/internal"
)

func makeReceiverContent() *fyne.Container {
	titleText := widget.NewLabel("> Receiver:")

	generateKeyPairButton := widget.NewButton("Generate Key Pair", func() {
		internal.GenerateKeypair()
		Log("Generate Key Pair.")
	})

	var hasPrivKey, hasMessage bool

	decryptCheckText := widget.NewMultiLineEntry()

	go func() {
		for range time.NewTicker(time.Second).C {
			fyne.Do(func() {
				hasPrivKey, hasMessage = isFileExist("priv.key", "CIPHER")

				text := fmt.Sprintf("Check 1 - Has Private Key. ('./priv.key' file): %t\n"+
					"Check 2 - Has Cipher File. ('./CIPHER.XXX' file): %t", hasPrivKey, hasMessage)
				decryptCheckText.SetText(text)
			})
		}
	}()

	decryptButton := widget.NewButton("Decrypt", func() {
		if !hasPrivKey || !hasMessage {
			Log("Not Ready for Decrypt.")
			return
		}

		isSuccess := internal.Decrypt()
		Log("Decrypt Message " + boolToString(isSuccess) + ".")
	})

	titleWrapper := container.NewBorder(blank(40), blank(40), blank(20), nil, titleText)

	content := container.NewVBox(generateKeyPairButton, blank(40), decryptCheckText, decryptButton)

	return container.NewBorder(titleWrapper, nil, blank(60), blank(60), content)
}
