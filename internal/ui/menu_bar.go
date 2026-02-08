package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MenuBar represents the top menu bar
type MenuBar struct {
	appActionsHandler AppActionsHandler
	pinHandler       PinHandler
	pinned           bool
	databaseLocation string
	container        *fyne.Container
	pinBtn           *widget.Button
}

// NewMenuBar creates a new menu bar
func NewMenuBar(appActions AppActionsHandler, pinHandler PinHandler, dbLocation string) *MenuBar {
	return &MenuBar{
		appActionsHandler: appActions,
		pinHandler:       pinHandler,
		databaseLocation: dbLocation,
	}
}

// Build builds the menu bar UI
func (m *MenuBar) Build() *fyne.Container {
	importBtn := widget.NewButton("Import", func() {
		m.appActionsHandler.OnImportNote()
	})
	
	exportBtn := widget.NewButton("Export", func() {
		m.appActionsHandler.OnExportNote()
	})
	
	createBtn := widget.NewButtonWithIcon("New Note", theme.DocumentCreateIcon(), func() {
		m.appActionsHandler.OnCreateNote()
	})
	
	m.pinBtn = widget.NewButtonWithIcon("Pin Mode", theme.InfoIcon(), func() {
		m.pinned = !m.pinned
		m.pinHandler.OnPinNote(m.pinned)
		if m.pinned {
			m.pinBtn.SetText("Unpin")
			m.pinBtn.SetIcon(theme.CancelIcon())
		} else {
			m.pinBtn.SetText("Pin Mode")
			m.pinBtn.SetIcon(theme.InfoIcon())
		}
	})
	
	dbLabel := widget.NewLabel(fmt.Sprintf("DB: %s", m.databaseLocation))
	dbLabel.TextStyle = fyne.TextStyle{Italic: true}
	
	m.container = container.NewHBox(
		importBtn,
		exportBtn,
		widget.NewSeparator(),
		createBtn,
		m.pinBtn,
		widget.NewSeparator(),
		dbLabel,
	)
	
	return m.container
}

// SetPinned updates the pin button state
func (m *MenuBar) SetPinned(pinned bool) {
	m.pinned = pinned
	if m.pinBtn != nil {
		if pinned {
			m.pinBtn.SetText("Unpin")
			m.pinBtn.SetIcon(theme.CancelIcon())
		} else {
			m.pinBtn.SetText("Pin Mode")
			m.pinBtn.SetIcon(theme.InfoIcon())
		}
	}
}