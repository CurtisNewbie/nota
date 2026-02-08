package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/curtisnewbie/nota/internal/domain"
)

// NoteEditor represents the note editor panel
type NoteEditor struct {
	editHandler   NoteEditHandler
	note          *domain.Note
	isSaving      bool
	minimalMode   bool
	titleEntry    *widget.Entry
	contentEntry  *widget.Entry
	createdLabel  *widget.Label
	updatedLabel  *widget.Label
	statusLabel   *widget.Label
	saveBtn       *widget.Button
	topBar        *fyne.Container
	bottomBar     *fyne.Container
	leftPanel     *fyne.Container
	container     *fyne.Container
}

// NewNoteEditor creates a new note editor
func NewNoteEditor(editHandler NoteEditHandler) *NoteEditor {
	return &NoteEditor{
		editHandler: editHandler,
	}
}

// Build builds the note editor UI
func (e *NoteEditor) Build() *fyne.Container {
	e.titleEntry = widget.NewEntry()
	e.titleEntry.SetPlaceHolder("Note Title")
	e.titleEntry.OnChanged = func(string) {
		if e.editHandler != nil && !e.isSaving {
			e.editHandler.OnContentChanged()
		}
	}

	e.contentEntry = widget.NewMultiLineEntry()
	e.contentEntry.SetPlaceHolder("Note content...")
	e.contentEntry.SetMinRowsVisible(20) // Increase default visible rows
	e.contentEntry.OnChanged = func(string) {
		if e.editHandler != nil && !e.isSaving {
			e.editHandler.OnContentChanged()
		}
	}

	e.createdLabel = widget.NewLabel("")
	e.createdLabel.TextStyle = fyne.TextStyle{Italic: true}

	e.updatedLabel = widget.NewLabel("")
	e.updatedLabel.TextStyle = fyne.TextStyle{Italic: true}

	e.statusLabel = widget.NewLabel("")
	e.statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	e.saveBtn = widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		if e.editHandler != nil {
			e.editHandler.OnSave()
		}
	})

	e.topBar = container.NewBorder(nil, nil, nil, e.saveBtn)

	e.bottomBar = container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(e.createdLabel, e.updatedLabel),
		e.statusLabel,
	)

	e.leftPanel = container.NewBorder(
		e.topBar,
		e.bottomBar,
		nil,
		nil,
		container.NewVBox(
			e.titleEntry,
			widget.NewSeparator(),
			e.contentEntry,
		),
	)

	e.container = container.NewBorder(nil, nil, nil, nil, e.leftPanel)

	return e.container
}

// DisplayNote displays a note in the editor
func (e *NoteEditor) DisplayNote(note *domain.Note) {
	e.note = note

	if note == nil {
		e.titleEntry.SetText("")
		e.contentEntry.SetText("")
		e.createdLabel.SetText("")
		e.updatedLabel.SetText("")
		e.statusLabel.SetText("No note selected")
		return
	}

	e.titleEntry.SetText(note.Title)
	e.contentEntry.SetText(note.Content)
	e.createdLabel.SetText(fmt.Sprintf("Created: %s", note.CreatedAt.Format("2006/01/02 15:04")))
	e.updatedLabel.SetText(fmt.Sprintf("Updated: %s", note.UpdatedAt.Format("2006/01/02 15:04")))
	e.statusLabel.SetText("Saved")

	// Don't set edit mode here - let the caller control it
	// e.setEditMode(false)
}

// GetTitle returns the current title
func (e *NoteEditor) GetTitle() string {
	return e.titleEntry.Text
}

// GetContent returns the current content
func (e *NoteEditor) GetContent() string {
	return e.contentEntry.Text
}

// MarkAsSaved marks the note as saved
func (e *NoteEditor) MarkAsSaved() {
	e.statusLabel.SetText("Saved")
	e.statusLabel.Importance = widget.LowImportance
}

// MarkAsUnsaved marks the note as unsaved
func (e *NoteEditor) MarkAsUnsaved() {
	e.statusLabel.SetText("Unsaved changes")
	e.statusLabel.Importance = widget.HighImportance
}

// StartSaving marks the start of a save operation
func (e *NoteEditor) StartSaving() {
	e.isSaving = true
}

// EndSaving marks the end of a save operation
func (e *NoteEditor) EndSaving() {
	e.isSaving = false
}

// ShowEmptyState shows the empty state
func (e *NoteEditor) ShowEmptyState() {
	e.titleEntry.SetText("")
	e.contentEntry.SetText("")
	e.createdLabel.SetText("")
	e.updatedLabel.SetText("")
	e.statusLabel.SetText("No notes available. Click 'New Note' to create one.")
}

// SetMinimalMode toggles minimal mode (hides UI elements)
func (e *NoteEditor) SetMinimalMode(minimal bool) {
	e.minimalMode = minimal
}
