package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

	n.noteList = widget.NewList(
		func() int {
			return len(n.notes)
		},
		func() fyne.CanvasObject {
			titleLabel := widget.NewLabel("")
			titleLabel.TextStyle = fyne.TextStyle{Bold: true}
			dateLabel := widget.NewLabel("")
			dateLabel.TextStyle = fyne.TextStyle{Italic: true}
			dateLabel.Importance = widget.LowImportance
			return container.NewVBox(titleLabel, dateLabel)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			container := obj.(*fyne.Container)
			if id >= 0 && id < len(n.notes) {
				note := n.notes[id]
				updatedTime := note.UpdatedAt.Format("2006/01/02")
				objects := container.Objects
				if len(objects) >= 2 {
					titleLabel := objects[0].(*widget.Label)
					dateLabel := objects[1].(*widget.Label)
					titleLabel.SetText(note.Title)
					dateLabel.SetText(updatedTime)
				}
			}
		},
	)

	n.noteList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(n.notes) && n.selectionHandler != nil {
			n.selectionHandler.OnNoteSelected(n.notes[id])
		}
	}

	rightPanel := container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		n.noteList,
	)

	n.rightPanel = rightPanel
	n.container = container.NewBorder(nil, nil, nil, nil, rightPanel)

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
	if n.window == nil {
		return
	}

	dialog.ShowConfirm("Delete Note",
		"Are you sure you want to delete this note?",
		func(confirmed bool) {
			if confirmed && n.deleteHandler != nil {
				n.deleteHandler.OnDeleteNote()
			}
		},
		n.window,
	)
}

// DisplayNotes displays the list of notes
func (n *NoteList) DisplayNotes(notes []*domain.Note) {
	n.notes = notes
	n.noteList.Refresh()
}

// RefreshNoteList refreshes the note list from the service
func (n *NoteList) RefreshNoteList() {
	n.noteList.Refresh()
}

// ClearSearch clears the search field
func (n *NoteList) ClearSearch() {
	n.searchEntry.SetText("")
}

// GetSearchQuery returns the current search query
func (n *NoteList) GetSearchQuery() string {
	return n.searchEntry.Text
}
