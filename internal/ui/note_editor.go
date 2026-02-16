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
	editHandler           NoteEditHandler
	deleteHandler         DeleteHandler
	note                  *domain.Note
	isSaving              bool
	minimalMode           bool
	titleEntry            *widget.Entry
	contentEntry          *widget.Entry
	createdLabel          *widget.Label
	updatedLabel          *widget.Label
	statusLabel           *widget.Label
	saveBtn               *widget.Button
	deleteBtn             *widget.Button
	topBar                *fyne.Container
	bottomBar             *fyne.Container
	leftPanel             *fyne.Container
	container             *fyne.Container
	minimizedTitleEntry   *widget.Entry
	minimizedContentEntry *widget.Entry
	minimizedStatusLabel  *widget.Label
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

	// Add shortcut to title entry's canvas if window is available
	// Note: This won't work when entry has focus, so we need a different approach
	// For now, this provides a fallback when entry doesn't have focus

	e.contentEntry = widget.NewMultiLineEntry()
	e.contentEntry.SetPlaceHolder("Note content...")
	e.contentEntry.SetMinRowsVisible(30) // Increase default visible rows
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

	e.saveBtn = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		if e.editHandler != nil {
			e.editHandler.OnSave()
		}
	})

	e.deleteBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if e.deleteHandler != nil {
			e.deleteHandler.OnDeleteNote()
		}
	})
	e.deleteBtn.Importance = widget.DangerImportance

	// Create button row with save and delete buttons
	buttonRow := container.NewHBox(e.saveBtn, e.deleteBtn)

	e.topBar = container.NewBorder(nil, nil, nil, buttonRow)

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
		// Disable fields when no note is selected
		e.titleEntry.Disable()
		e.contentEntry.Disable()
		return
	}

	e.titleEntry.SetText(note.Title)
	e.contentEntry.SetText(note.Content)
	e.createdLabel.SetText(fmt.Sprintf("Created: %s", note.CreatedAt.Format("2006/01/02 15:04")))
	e.updatedLabel.SetText(fmt.Sprintf("Updated: %s", note.UpdatedAt.Format("2006/01/02 15:04")))
	e.statusLabel.SetText("Saved")
	// Enable fields when a note is displayed
	e.titleEntry.Enable()
	e.contentEntry.Enable()

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
	// Update minimized status label if it exists
	if e.minimizedStatusLabel != nil {
		e.minimizedStatusLabel.SetText("Saved")
		e.minimizedStatusLabel.Importance = widget.LowImportance
	}
}

// MarkAsUnsaved marks the note as unsaved
func (e *NoteEditor) MarkAsUnsaved() {
	e.statusLabel.SetText("Unsaved changes")
	e.statusLabel.Importance = widget.HighImportance
	// Update minimized status label if it exists
	if e.minimizedStatusLabel != nil {
		e.minimizedStatusLabel.SetText("Unsaved changes")
		e.minimizedStatusLabel.Importance = widget.HighImportance
	}
}

// SetMinimizedStatusLabel sets the status label for minimized mode
func (e *NoteEditor) SetMinimizedStatusLabel(statusLabel *widget.Label) {
	e.minimizedStatusLabel = statusLabel
	// Initialize with current status
	if e.minimizedStatusLabel != nil && e.statusLabel != nil {
		e.minimizedStatusLabel.SetText(e.statusLabel.Text)
		e.minimizedStatusLabel.Importance = e.statusLabel.Importance
	}
}

// StartSaving marks the start of a save operation
func (e *NoteEditor) StartSaving() {
	e.isSaving = true
}

// EndSaving marks the end of a save operation
func (e *NoteEditor) EndSaving() {
	e.isSaving = false
}

// IsSaving returns whether a save operation is in progress
func (e *NoteEditor) IsSaving() bool {
	return e.isSaving
}

// SetDeleteHandler sets the delete handler for the note editor
func (e *NoteEditor) SetDeleteHandler(handler DeleteHandler) {
	e.deleteHandler = handler
}

// ShowEmptyState shows the empty state
func (e *NoteEditor) ShowEmptyState() {
	e.titleEntry.SetText("")
	e.contentEntry.SetText("")
	e.createdLabel.SetText("")
	e.updatedLabel.SetText("")
	e.statusLabel.SetText("No notes available. Click 'New Note' to create one.")
	// Disable fields when no notes are available
	e.titleEntry.Disable()
	e.contentEntry.Disable()
}

// SetMinimalMode toggles minimal mode (hides UI elements)
func (e *NoteEditor) SetMinimalMode(minimal bool) {
	e.minimalMode = minimal
}

// SetMinimizedWidgets stores references to the minimized mode widgets
func (e *NoteEditor) SetMinimizedWidgets(titleEntry *widget.Entry, contentEntry *widget.Entry) {
	e.minimizedTitleEntry = titleEntry
	e.minimizedContentEntry = contentEntry
}

// SyncFromMinimizedMode syncs changes from minimized mode widgets back to the main editor
func (e *NoteEditor) SyncFromMinimizedMode() {
	if e.minimizedTitleEntry != nil && e.minimizedContentEntry != nil {
		// Get the current values from minimized widgets
		newTitle := e.minimizedTitleEntry.Text
		newContent := e.minimizedContentEntry.Text

		// Get the current values from main editor
		currentTitle := e.titleEntry.Text
		currentContent := e.contentEntry.Text

		// Only update and mark as unsaved if values actually changed
		if newTitle != currentTitle || newContent != currentContent {
			e.titleEntry.SetText(newTitle)
			e.contentEntry.SetText(newContent)
			// Mark as unsaved since content changed
			if e.editHandler != nil {
				e.editHandler.OnContentChanged()
			}
		}

		// Clear references
		e.minimizedTitleEntry = nil
		e.minimizedContentEntry = nil
	}
}
