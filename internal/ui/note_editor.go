package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/curtisnewbie/nota/internal/domain"
)

// NoteEditor represents the note editor/viewer panel
type NoteEditor struct {
	editHandler  NoteEditHandler
	mainUI       *MainUI
	note         *domain.Note
	isEditMode   bool
	isSaving     bool
	titleEntry   *widget.Entry
	contentEntry *widget.Entry
	createdLabel *widget.Label
	updatedLabel *widget.Label
	statusLabel  *widget.Label
	modeBtn      *widget.Button
	container    *fyne.Container
}

// NewNoteEditor creates a new note editor
func NewNoteEditor(editHandler NoteEditHandler) *NoteEditor {
	return &NoteEditor{
		editHandler: editHandler,
	}
}

// SetMainUI sets the main UI reference
func (e *NoteEditor) SetMainUI(mainUI *MainUI) {
	e.mainUI = mainUI
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

	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		if e.editHandler != nil {
			e.editHandler.OnSave()
		}
	})

	e.modeBtn = widget.NewButton("Edit", func() {
		e.toggleEditMode()
	})

	topBar := container.NewBorder(nil, nil, nil, container.NewHBox(e.modeBtn, saveBtn))

	bottomBar := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(e.createdLabel, e.updatedLabel),
		e.statusLabel,
	)

	leftPanel := container.NewBorder(
		topBar,
		bottomBar,
		nil,
		nil,
		container.NewVBox(
			e.titleEntry,
			widget.NewSeparator(),
			e.contentEntry,
		),
	)

	e.container = container.NewBorder(nil, nil, nil, nil, leftPanel)

	e.setEditMode(false)

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
		e.setEditMode(false)
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
	e.setEditMode(false)
}

// toggleEditMode toggles between edit and read mode
func (e *NoteEditor) toggleEditMode() {
	e.setEditMode(!e.isEditMode)
}

// setEditMode sets the edit mode
func (e *NoteEditor) setEditMode(editMode bool) {
	e.isEditMode = editMode

	if editMode {
		e.titleEntry.Enable()
		e.contentEntry.Enable()
		if e.modeBtn != nil {
			e.modeBtn.SetText("View")
		}
	} else {
		// Disable entries in view mode to make them read-only
		e.titleEntry.Disable()
		e.contentEntry.Disable()
		if e.modeBtn != nil {
			e.modeBtn.SetText("Edit")
		}
	}

	// Don't call mainUI.SetEditMode here to avoid circular call
	// MainUI.SetEditMode now calls EnableEdit/DisableEdit directly
}

// EnableEdit enables edit mode
func (e *NoteEditor) EnableEdit() {
	e.setEditMode(true)
}

// DisableEdit disables edit mode (read-only)
func (e *NoteEditor) DisableEdit() {
	e.setEditMode(false)
}

// IsEditMode returns whether currently in edit mode
func (e *NoteEditor) IsEditMode() bool {
	return e.isEditMode
}
