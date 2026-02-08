package ui

import (

	"fmt"



	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"

	"github.com/curtisnewbie/nota/internal/domain"

)

// NoteList represents the note list panel
type NoteList struct {
	selectionHandler NoteSelectionHandler
	searchHandler    SearchHandler
	notes            []*domain.Note
	searchEntry      *widget.Entry
	noteList         *widget.List
	container        *fyne.Container
}

// NewNoteList creates a new note list
func NewNoteList(selectionHandler NoteSelectionHandler, searchHandler SearchHandler) *NoteList {
	return &NoteList{
		selectionHandler: selectionHandler,
		searchHandler:    searchHandler,
	}
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
	
	n.noteList = widget.NewList(
		func() int {
			return len(n.notes)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id >= 0 && id < len(n.notes) {
				note := n.notes[id]
				updatedTime := note.UpdatedAt.Format("2006/01/02")
				label.SetText(fmt.Sprintf("%s\n%s", note.Title, updatedTime))
			}
		},
	)
	
	n.noteList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(n.notes) && n.selectionHandler != nil {
			n.selectionHandler.OnNoteSelected(n.notes[id])
		}
	}
	
	rightPanel := container.NewBorder(
		n.searchEntry,
		nil,
		nil,
		nil,
		n.noteList,
	)
	
	n.container = container.NewBorder(nil, nil, nil, nil, rightPanel)
	
	return n.container
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
