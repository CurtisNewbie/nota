package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	window      fyne.Window
	app         AppActionsHandler
	menuBar     *MenuBar
	noteEditor  *NoteEditor
	noteList    *NoteList
	container   *fyne.Container
	noteService NoteService
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
	menuBarContainer := m.menuBar.Build()

	leftPanel := m.noteList.Build()
	rightPanel := m.noteEditor.Build()

	splitContainer := container.NewHSplit(leftPanel, rightPanel)
	splitContainer.SetOffset(0.4)

	mainContainer := container.NewBorder(
		menuBarContainer,
		nil,
		nil,
		nil,
		splitContainer,
	)

	m.container = container.NewWithoutLayout(mainContainer)

	// Add right-click detection for the note list
	m.setupRightClickHandler()

	return mainContainer
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
