package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/curtisnewbie/nota/internal/i18n"
)

// MenuBar represents the top menu bar
type MenuBar struct {
	appActionsHandler AppActionsHandler
	pinHandler        PinHandler
	languageHandler   LanguageHandler
	pinned            bool
	databaseLocation  string
	container         *fyne.Container
	window            fyne.Window
}

// NewMenuBar creates a new menu bar
func NewMenuBar(appActions AppActionsHandler, pinHandler PinHandler, languageHandler LanguageHandler, dbLocation string) *MenuBar {
	return &MenuBar{
		appActionsHandler: appActions,
		pinHandler:        pinHandler,
		languageHandler:   languageHandler,
		databaseLocation:  dbLocation,
	}
}

// SetWindow sets the window reference for popup menus
func (m *MenuBar) SetWindow(window fyne.Window) {
	m.window = window
}

// Build builds the menu bar UI
func (m *MenuBar) Build() *fyne.Container {
	t := i18n.T()

	noteBtn := widget.NewButton(t.Menu.Note, func() {
		m.showNoteMenu()
	})

	fileBtn := widget.NewButton(t.Menu.File, func() {
		m.showFileMenu()
	})

	viewBtn := widget.NewButton(t.Menu.View, func() {
		m.showViewMenu()
	})

	languageBtn := widget.NewButton(t.Menu.Language, func() {
		m.showLanguageMenu()
	})

	dbLabel := widget.NewLabel(fmt.Sprintf(t.Database.Location, m.databaseLocation))
	dbLabel.TextStyle = fyne.TextStyle{Italic: true}

	m.container = container.NewHBox(
		noteBtn,
		fileBtn,
		viewBtn,
		languageBtn,
		widget.NewSeparator(),
		dbLabel,
	)

	return m.container
}

// showNoteMenu shows the Note dropdown menu
func (m *MenuBar) showNoteMenu() {
	if m.window == nil {
		return
	}

	t := i18n.T()

	menu := fyne.NewMenu("",
		fyne.NewMenuItem(t.Menu.NewNote, func() {
			m.appActionsHandler.OnCreateNote()
		}),
	)

	popUp := widget.NewPopUpMenu(menu, m.window.Canvas())
	pos := m.menuButtonPosition(0)
	popUp.Move(pos)
	popUp.Show()
}

// showFileMenu shows the File dropdown menu
func (m *MenuBar) showFileMenu() {
	if m.window == nil {
		return
	}

	t := i18n.T()

	menu := fyne.NewMenu("",
		fyne.NewMenuItem(t.Menu.Import, func() {
			m.appActionsHandler.OnImportNote()
		}),
		fyne.NewMenuItem(t.Menu.Export, func() {
			m.appActionsHandler.OnExportNote()
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

	t := i18n.T()

	menu := fyne.NewMenu("",
		fyne.NewMenuItem(t.Menu.MinimizedMode, func() {
			m.togglePinMode()
		}),
	)

	popUp := widget.NewPopUpMenu(menu, m.window.Canvas())
	pos := m.menuButtonPosition(2)
	popUp.Move(pos)
	popUp.Show()
}

// showLanguageMenu shows the Language dropdown menu
func (m *MenuBar) showLanguageMenu() {
	if m.window == nil {
		return
	}

	t := i18n.T()

	menu := fyne.NewMenu("",
		fyne.NewMenuItem(t.Menu.English, func() {
			if m.languageHandler != nil {
				m.languageHandler.OnLanguageChanged(i18n.LanguageEnglish)
			}
		}),
		fyne.NewMenuItem(t.Menu.Chinese, func() {
			if m.languageHandler != nil {
				m.languageHandler.OnLanguageChanged(i18n.LanguageChinese)
			}
		}),
	)

	popUp := widget.NewPopUpMenu(menu, m.window.Canvas())
	pos := m.menuButtonPosition(3)
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

// Refresh refreshes the menu bar UI (used when language changes)
func (m *MenuBar) Refresh() {
	t := i18n.T()
	buttons := m.container.Objects

	// Update button texts
	if len(buttons) >= 4 {
		if noteBtn, ok := buttons[0].(*widget.Button); ok {
			noteBtn.SetText(t.Menu.Note)
		}
		if fileBtn, ok := buttons[1].(*widget.Button); ok {
			fileBtn.SetText(t.Menu.File)
		}
		if viewBtn, ok := buttons[2].(*widget.Button); ok {
			viewBtn.SetText(t.Menu.View)
		}
		if langBtn, ok := buttons[3].(*widget.Button); ok {
			langBtn.SetText(t.Menu.Language)
		}
	}

	// Update database label
	if len(buttons) >= 5 {
		if dbLabel, ok := buttons[5].(*widget.Label); ok {
			dbLabel.SetText(fmt.Sprintf(t.Database.Location, m.databaseLocation))
		}
	}
}
