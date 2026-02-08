package ui

import (
	"github.com/curtisnewbie/nota/internal/domain"
)

// NoteSelectionHandler handles note selection events
type NoteSelectionHandler interface {
	OnNoteSelected(note *domain.Note)
}

// NoteEditHandler handles note edit events
type NoteEditHandler interface {
	OnContentChanged()
	OnSave()
}

// AppActionsHandler handles application action events
type AppActionsHandler interface {
	OnCreateNote()
	OnDeleteNote()
	OnImportNote()
	OnExportNote()
	OnNoteSelected(note *domain.Note)
	OnContentChanged()
	OnSave()
	OnSearch(query string)
	OnPinNote(pin bool)
	GetDatabaseLocation() string
	ListNotes() ([]*domain.Note, error)
}

// SearchHandler handles search events
type SearchHandler interface {
	OnSearch(query string)
}

// PinHandler handles pin mode events
type PinHandler interface {
	OnPinNote(pin bool)
}
