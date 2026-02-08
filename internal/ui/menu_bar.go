package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// MenuBar represents the top menu bar
type MenuBar struct {
	appActionsHandler AppActionsHandler
	pinHandler        PinHandler
	pinned            bool
	databaseLocation  string
	container         *fyne.Container
	window            fyne.Window
	editMode          bool
	editBtn           *widget.Button
}

// NewMenuBar creates a new menu bar
func NewMenuBar(appActions AppActionsHandler, pinHandler PinHandler, dbLocation string) *MenuBar {
	return &MenuBar{
		appActionsHandler: appActions,
		pinHandler:        pinHandler,
		databaseLocation:  dbLocation,
		editMode:          false,
	}
}

// SetWindow sets the window reference for popup menus
func (m *MenuBar) SetWindow(window fyne.Window) {
	m.window = window
}

// Build builds the menu bar UI
func (m *MenuBar) Build() *fyne.Container {
	fileBtn := widget.NewButton("File", func() {
		m.showFileMenu()
	})

	editBtn := widget.NewButton("Edit Mode", func() {
		m.showEditMenu()
	})
	m.editBtn = editBtn

	viewBtn := widget.NewButton("View", func() {
		m.showViewMenu()
	})

	dbLabel := widget.NewLabel(fmt.Sprintf("DB: %s", m.databaseLocation))
	dbLabel.TextStyle = fyne.TextStyle{Italic: true}

	m.container = container.NewHBox(
		fileBtn,
		editBtn,
		viewBtn,
		widget.NewSeparator(),
		dbLabel,
	)

	return m.container
}

// showFileMenu shows the File dropdown menu
func (m *MenuBar) showFileMenu() {
	if m.window == nil {
		return
	}

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("New Note", func() {
			m.appActionsHandler.OnCreateNote()
		}),
		fyne.NewMenuItem("Import", func() {
			m.appActionsHandler.OnImportNote()
		}),
		fyne.NewMenuItem("Export", func() {
			m.appActionsHandler.OnExportNote()
		}),
	)

	popUp := widget.NewPopUpMenu(menu, m.window.Canvas())
	pos := m.menuButtonPosition(0)
	popUp.Move(pos)
	popUp.Show()
}

// showEditMenu shows the Edit dropdown menu
func (m *MenuBar) showEditMenu() {
	if m.window == nil {
		return
	}

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Delete Note", func() {
			m.appActionsHandler.OnDeleteNote()
		}),
	)

	popUp := widget.NewPopUpMenu(menu, m.window.Canvas())
	pos := m.menuButtonPosition(1)
	popUp.Move(pos)
	popUp.Show()
}

// showViewMenu shows the View dropdown menu
func (m *MenuBar) showViewMenu() {
	if m.window == nil {
		return
	}

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Pin Mode", func() {
			m.togglePinMode()
		}),
	)

	popUp := widget.NewPopUpMenu(menu, m.window.Canvas())
	pos := m.menuButtonPosition(2)
	popUp.Move(pos)
	popUp.Show()
}

// menuButtonPosition calculates the position for a menu popup based on button index
func (m *MenuBar) menuButtonPosition(index int) fyne.Position {
	if m.window == nil {
		return fyne.NewPos(10, 30)
	}

	pos := m.container.Position()
	buttonWidth := 80           // approximate width per button
	buttonHeight := float32(40) // approximate height

	return fyne.NewPos(pos.X+float32(index*buttonWidth), pos.Y+buttonHeight)
}

// togglePinMode toggles pin mode
func (m *MenuBar) togglePinMode() {
	m.pinned = !m.pinned
	m.pinHandler.OnPinNote(m.pinned)
}

// SetPinned updates the pin mode state
func (m *MenuBar) SetPinned(pinned bool) {
	m.pinned = pinned
}

// SetEditMode updates the edit mode state and updates the Edit Mode button text
func (m *MenuBar) SetEditMode(editMode bool) {
	m.editMode = editMode
	if m.editBtn != nil {
		if editMode {
			m.editBtn.SetText("View Mode")
		} else {
			m.editBtn.SetText("Edit Mode")
		}
	}
}
