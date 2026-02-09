package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/nota/internal/domain"
)

// NoteService defines the interface for note operations
type NoteService interface {
	ListNotes(rail flow.Rail) ([]*domain.Note, error)
}

// ImportExportService defines the interface for import/export operations
type ImportExportService interface{}

// MainUI represents the main UI
type MainUI struct {
	window           fyne.Window
	app              AppActionsHandler
	menuBar          *MenuBar
	noteEditor       *NoteEditor
	noteList         *NoteList
	container        *fyne.Container
	noteService      NoteService
	minimized        bool
	menuBarContainer *fyne.Container
	rightPanel       *fyne.Container
	fullContainer    *fyne.Container
}

// NewMainUI creates a new main UI
func NewMainUI(
	window fyne.Window,
	noteService NoteService,
	importExportService ImportExportService,
	app AppActionsHandler,
) *MainUI {
	mainUI := &MainUI{
		window:      window,
		app:         app,
		noteService: noteService,
	}

	mainUI.menuBar = NewMenuBar(app, app, app.GetDatabaseLocation())
	mainUI.menuBar.SetWindow(window)
	mainUI.noteEditor = NewNoteEditor(app)
	mainUI.noteList = NewNoteList(app, app)
	mainUI.noteList.SetWindow(window)
	mainUI.noteList.SetDeleteHandler(app)

	return mainUI
}

// Build builds the main UI
func (m *MainUI) Build() fyne.CanvasObject {
	m.menuBarContainer = m.menuBar.Build()

	leftPanel := m.noteList.Build()
	m.rightPanel = m.noteEditor.Build()

	splitContainer := container.NewHSplit(leftPanel, m.rightPanel)
	splitContainer.SetOffset(0.17)

	m.fullContainer = container.NewBorder(
		m.menuBarContainer,
		nil,
		nil,
		nil,
		splitContainer,
	)

	m.container = m.fullContainer

	// Add right-click detection for the note list
	m.setupRightClickHandler()

	return m.container
}

// setupRightClickHandler sets up right-click detection for the note list
func (m *MainUI) setupRightClickHandler() {
	if m.window == nil || m.noteList == nil {
		return
	}

	// Right-click handling will be implemented differently
	// This is a placeholder for future implementation
}

// DisplayNote displays a note
func (m *MainUI) DisplayNote(note *domain.Note) {
	m.noteEditor.DisplayNote(note)
}

// ShowEmptyState shows the empty state
func (m *MainUI) ShowEmptyState() {
	m.noteEditor.ShowEmptyState()
	m.noteList.DisplayNotes([]*domain.Note{})
}

// RefreshNoteList refreshes the note list
func (m *MainUI) RefreshNoteList() {
	notes, err := m.noteService.ListNotes(flow.NewRail(context.Background()))
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	m.noteList.DisplayNotes(notes)
}

// DisplaySearchResults displays search results
func (m *MainUI) DisplaySearchResults(notes []*domain.Note) {
	m.noteList.DisplayNotes(notes)
}

// GetTitle returns the current title
func (m *MainUI) GetTitle() string {
	return m.noteEditor.GetTitle()
}

// GetContent returns the current content
func (m *MainUI) GetContent() string {
	return m.noteEditor.GetContent()
}

// MarkAsSaved marks the note as saved
func (m *MainUI) MarkAsSaved() {
	m.noteEditor.MarkAsSaved()
}

// MarkAsUnsaved marks the note as unsaved
func (m *MainUI) MarkAsUnsaved() {
	m.noteEditor.MarkAsUnsaved()
}

// SetPinned sets the pin mode
func (m *MainUI) SetPinned(pinned bool) {
	m.menuBar.SetPinned(pinned)
}

// StartSaving marks the start of a save operation
func (m *MainUI) StartSaving() {
	m.noteEditor.StartSaving()
}

// EndSaving marks the end of a save operation
func (m *MainUI) EndSaving() {
	m.noteEditor.EndSaving()
}

// ToggleMinimizedMode toggles between normal and minimized (notepad) mode
func (m *MainUI) ToggleMinimizedMode(minimized bool) {
	m.minimized = minimized

	if minimized {
		// Create minimal container with title, content, and exit button
		exitBtn := widget.NewButton("Exit Minimized Mode", func() {
			m.ToggleMinimizedMode(false)
			m.SetPinned(false)
			m.window.Resize(fyne.NewSize(1200, 800))
		})
		exitBtn.Importance = widget.MediumImportance

		// Create fresh widget instances for minimized mode to avoid state conflicts
		titleEntry := widget.NewEntry()
		titleEntry.SetText(m.noteEditor.GetTitle())
		titleEntry.PlaceHolder = "Note Title"

		contentEntry := widget.NewMultiLineEntry()
		contentEntry.SetText(m.noteEditor.GetContent())
		contentEntry.SetPlaceHolder("Note Content")
		contentEntry.Wrapping = fyne.TextWrapWord
		contentEntry.SetMinRowsVisible(20)

		minimalContainer := container.NewVBox(
			exitBtn,
			titleEntry,
			widget.NewSeparator(),
			contentEntry,
		)
		m.container = minimalContainer
	} else {
		// Restore full container
		m.container = m.fullContainer
	}

	// Update window content
	m.window.SetContent(m.container)
}
