package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme struct {
	fyne.Theme
}

func (t *Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) (c color.Color) {
	switch name {
	case theme.ColorNameBackground:
		c = color_MainDark
	case theme.ColorNameButton:
		c = color_MainLight
	case theme.ColorNameInputBackground:
		c = color_MainLight
	default:
		c = theme.DefaultTheme().Color(name, variant)
	}

	return
}
