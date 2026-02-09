package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/curtisnewbie/nota/internal/domain"
)

// DeleteHandler handles note deletion events
type DeleteHandler interface {
	OnDeleteNote()
}

// NoteList represents the note list panel
type NoteList struct {
	selectionHandler NoteSelectionHandler
	searchHandler    SearchHandler
	deleteHandler    DeleteHandler
	window           fyne.Window
	notes            []*domain.Note
	searchEntry      *widget.Entry
	noteList         *widget.List
	rightPanel       *fyne.Container
	container        *fyne.Container
	menu             *fyne.Menu
	popUpMenu        *widget.PopUpMenu
	notesContainer   *fyne.Container
}

// NewNoteList creates a new note list
func NewNoteList(selectionHandler NoteSelectionHandler, searchHandler SearchHandler) *NoteList {
	return &NoteList{
		selectionHandler: selectionHandler,
		searchHandler:    searchHandler,
	}
}

// SetWindow sets the window for the note list (needed for dialogs)
func (n *NoteList) SetWindow(window fyne.Window) {
	n.window = window
}

// SetDeleteHandler sets the delete handler for the note list
func (n *NoteList) SetDeleteHandler(handler DeleteHandler) {
	n.deleteHandler = handler
}

// Build builds the note list UI
func (n *NoteList) Build() *fyne.Container {
	n.searchEntry = widget.NewEntry()
	n.searchEntry.SetPlaceHolder("Search notes...")
	n.searchEntry.OnChanged = func(query string) {
		if n.searchHandler != nil {
			n.searchHandler.OnSearch(query)
		}
	}

	// Add delete button
	deleteBtn := widget.NewButton("Delete", func() {
		n.onDeleteRequested()
	})
	deleteBtn.Importance = widget.DangerImportance

	// Create toolbar with search and delete
	// Put search entry in center (expandable) and delete button on right edge
	toolbar := container.NewBorder(
		nil,
		nil,
		nil,
		deleteBtn,
		n.searchEntry,
	)

	// Use a dummy list widget - actual notes will be displayed in notesContainer
	n.noteList = widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return container.NewVBox() },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)

	n.noteList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(n.notes) && n.selectionHandler != nil {
			n.selectionHandler.OnNoteSelected(n.notes[id])
		}
	}

	// Create notesContainer for individual note buttons
	n.notesContainer = container.NewVBox()
	notesScroll := container.NewScroll(n.notesContainer)
	notesScroll.SetMinSize(fyne.NewSize(150, 400))

	// Store the notes container for updates
	n.rightPanel = container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		notesScroll,
	)

	n.container = container.NewBorder(nil, nil, nil, nil, n.rightPanel)

	return n.container
}

// ShowContextMenu shows the context menu at the given position
func (n *NoteList) ShowContextMenu(pos fyne.Position) {
	if n.window == nil {
		return
	}

	// Create popup menu
	deleteItem := fyne.NewMenuItem("Delete", func() {
		n.onDeleteRequested()
	})

	menu := fyne.NewMenu("", deleteItem)
	n.popUpMenu = widget.NewPopUpMenu(menu, n.window.Canvas())
	n.popUpMenu.Move(pos)
}

// onDeleteRequested handles the delete request from context menu
func (n *NoteList) onDeleteRequested() {
	if n.deleteHandler != nil {
		n.deleteHandler.OnDeleteNote()
	}
}

// DisplayNotes displays the list of notes
func (n *NoteList) DisplayNotes(notes []*domain.Note) {
	n.notes = notes

	// Clear existing notes
	if n.notesContainer != nil {
		n.notesContainer.Objects = nil
		n.notesContainer.Refresh()
	}

	// Add note buttons
	for _, note := range notes {
		// Create a button that wraps the note content
		// Buttons always fire their OnTapped callback when clicked, even if already "selected"
		noteButton := widget.NewButton(note.Title+"   "+note.UpdatedAt.Format("2006/01/02"), func() {
			if n.selectionHandler != nil {
				n.selectionHandler.OnNoteSelected(note)
			}
		})
		noteButton.Alignment = widget.ButtonAlignLeading

		n.notesContainer.Add(noteButton)
	}

	if n.notesContainer != nil {
		n.notesContainer.Refresh()
	}
}

// RefreshNoteList refreshes the note list from the service
func (n *NoteList) RefreshNoteList() {
	if n.notesContainer != nil {
		n.notesContainer.Refresh()
	}
}

// ClearSearch clears the search field
func (n *NoteList) ClearSearch() {
	n.searchEntry.SetText("")
}

// GetSearchQuery returns the current search query
func (n *NoteList) GetSearchQuery() string {
	return n.searchEntry.Text
}
