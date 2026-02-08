package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// MaterialTheme implements fyne.Theme with Material Design colors
type MaterialTheme struct{}

var _ fyne.Theme = (*MaterialTheme)(nil)

func (t *MaterialTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.RGBA{R: 63, G: 81, B: 181, A: 255} // Material Indigo 500
	case theme.ColorNameBackground:
		if variant == theme.VariantDark {
			return color.RGBA{R: 32, G: 33, B: 36, A: 255} // Dark Gray
		}
		return color.RGBA{R: 250, G: 250, B: 250, A: 255} // Light Gray
	case theme.ColorNameForeground:
		if variant == theme.VariantDark {
			return color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		return color.RGBA{R: 33, G: 33, B: 33, A: 255}
	case theme.ColorNameButton:
		return color.RGBA{R: 63, G: 81, B: 181, A: 255}
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 189, G: 189, B: 189, A: 255}
	case theme.ColorNameHover:
		return color.RGBA{R: 63, G: 81, B: 181, A: 255}
	case theme.ColorNameFocus:
		return color.RGBA{R: 63, G: 81, B: 181, A: 255}
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 189, G: 189, B: 189, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 189, G: 189, B: 189, A: 255}
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 63, G: 81, B: 181, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 63, G: 81, B: 181, A: 100}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t *MaterialTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *MaterialTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *MaterialTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 12
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameInputBorder:
		return 1
	}
	return theme.DefaultTheme().Size(name)
}