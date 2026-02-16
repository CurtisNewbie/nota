package app

import (
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/nota/internal/domain"
	"github.com/curtisnewbie/nota/internal/infrastructure"
	"github.com/curtisnewbie/nota/internal/repository"
	"github.com/curtisnewbie/nota/internal/service"
	"github.com/curtisnewbie/nota/internal/ui"
)

// App represents the main application
type App struct {
	fyneApp             fyne.App
	window              fyne.Window
	noteService         service.NoteService
	importExportService service.ImportExportService
	mainUI              *ui.MainUI
	currentNote         *domain.Note
	hasUnsavedChanges   bool
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	rail := flow.EmptyRail()

	rail.Infof("Initializing Nota application...")

	fyneApp := app.New()
	fyneApp.Settings().SetTheme(&ui.MaterialTheme{})

	err := infrastructure.EnsureDatabaseDir()
	if err != nil {
		rail.Errorf("Failed to create database directory: %v", err)
		return nil, err
	}

	db, err := infrastructure.InitializeDatabase()
	if err != nil {
		rail.Errorf("Failed to initialize database: %v", err)
		return nil, err
	}

	noteRepo := repository.NewSQLiteNoteRepository(db)
	noteService := service.NewNoteService(noteRepo)
	importExportService := service.NewImportExportService(noteRepo)

	appInstance := &App{
		fyneApp:             fyneApp,
		noteService:         noteService,
		importExportService: importExportService,
	}

	window := fyneApp.NewWindow("Nota")
	window.Resize(fyne.NewSize(1200, 800))
	window.SetCloseIntercept(func() {
		appInstance.onClose()
	})

	appInstance.window = window

	mainUI := ui.NewMainUI(window, noteService, importExportService, appInstance)
	appInstance.mainUI = mainUI

	window.SetContent(mainUI.Build())

	// Refresh the note list on startup
	mainUI.RefreshNoteList()

	err = appInstance.loadLastNote()
	if err != nil {
		rail.Infof("No existing notes, ready to create new note")
		mainUI.ShowEmptyState()
	}

	rail.Infof("Application initialized successfully")

	return appInstance, nil
}

// Run starts the application
func (a *App) Run() {
	// Add keyboard shortcut for saving (Cmd+S on macOS, Ctrl+S on Windows/Linux)
	// Add before showing the window to ensure proper registration
	canvas := a.window.Canvas()
	canvas.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierShortcutDefault,
	}, func(shortcut fyne.Shortcut) {
		a.saveCurrentNote()
	})

	a.window.ShowAndRun()
}

// onClose handles window close event
func (a *App) onClose() {
	if a.hasUnsavedChanges && !a.isNoteEmpty() {
		dialog.ShowConfirm("Unsaved Changes",
			"You have unsaved changes. Do you want to save them before closing?",
			func(save bool) {
				if save {
					a.saveCurrentNote()
				}
				a.fyneApp.Quit()
			},
			a.window,
		)
	} else {
		a.fyneApp.Quit()
	}
}

// loadLastNote loads the last modified note
func (a *App) loadLastNote() error {
	rail := flow.EmptyRail()
	note, err := a.noteService.GetLastModifiedNote(rail)
	if err != nil {
		return err
	}
	a.mainUI.StartSaving()
	defer a.mainUI.EndSaving()
	a.currentNote = note
	a.mainUI.DisplayNote(note)
	a.mainUI.MarkAsSaved()
	return nil
}

// isNoteEmpty checks if the current note is a new, empty note (no content)
func (a *App) isNoteEmpty() bool {
	return a.currentNote != nil &&
		a.currentNote.ID == "" &&
		a.mainUI.GetTitle() == "" &&
		a.mainUI.GetContent() == ""
}

// saveCurrentNote saves the current note
func (a *App) saveCurrentNote() {
	if a.currentNote == nil {
		return
	}

	rail := flow.EmptyRail()
	a.mainUI.StartSaving()
	defer a.mainUI.EndSaving()

	a.currentNote.Title = a.mainUI.GetTitle()
	a.currentNote.Content = a.mainUI.GetContent()

	var err error
	isNewNote := a.currentNote.ID == ""

	if isNewNote {
		// Create new note in database
		err = a.noteService.CreateNote(rail, a.currentNote)
	} else {
		// Update existing note
		err = a.noteService.UpdateNote(rail, a.currentNote)
	}

	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	// Fetch latest note from database to get updated timestamps and other fields
	latestNote, fetchErr := a.noteService.GetNote(rail, a.currentNote.ID)
	if fetchErr != nil {
		dialog.ShowError(fetchErr, a.window)
		return
	}

	a.currentNote = latestNote
	a.mainUI.StartSaving()
	a.mainUI.DisplayNote(latestNote) // Update UI with latest note data
	a.mainUI.EndSaving()
	a.hasUnsavedChanges = false
	a.mainUI.MarkAsSaved()

	// Always refresh the note list after saving to show the latest changes
	a.mainUI.RefreshNoteList()
}

// onNoteSelected is called when a note is selected from the list
func (a *App) onNoteSelected(note *domain.Note) {
	if a.hasUnsavedChanges && !a.isNoteEmpty() {
		dialog.ShowConfirm("Unsaved Changes",
			"You have unsaved changes. Do you want to save them before switching?",
			func(save bool) {
				if save {
					a.saveCurrentNote()
				}
				// Always load the selected note as the current note, regardless of save choice
				rail := flow.EmptyRail()
				latestNote, err := a.noteService.GetNote(rail, note.ID)
				if err != nil {
					dialog.ShowError(err, a.window)
					return
				}
				a.mainUI.StartSaving()
				defer a.mainUI.EndSaving()
				a.currentNote = latestNote
				a.hasUnsavedChanges = false
				a.mainUI.DisplayNote(latestNote)
				a.mainUI.MarkAsSaved()
			},
			a.window,
		)
	} else {
		// Fetch latest note from database
		rail := flow.EmptyRail()
		latestNote, err := a.noteService.GetNote(rail, note.ID)
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		a.mainUI.StartSaving()
		defer a.mainUI.EndSaving()
		a.currentNote = latestNote
		a.hasUnsavedChanges = false
		a.mainUI.DisplayNote(latestNote)
		a.mainUI.MarkAsSaved()
	}
}

// onContentChanged is called when note content is modified
func (a *App) onContentChanged() {
	a.hasUnsavedChanges = true
	a.mainUI.MarkAsUnsaved()
}

// onCreateNote is called when user wants to create a new note
func (a *App) onCreateNote() {
	if a.hasUnsavedChanges && !a.isNoteEmpty() {
		dialog.ShowConfirm("Unsaved Changes",
			"You have unsaved changes. Do you want to save them before creating a new note?",
			func(save bool) {
				if save {
					a.saveCurrentNote()
				}
				a.createNewNote()
			},
			a.window,
		)
	} else {
		a.createNewNote()
	}
}

// createNewNote creates a new note in memory (not saved to database yet)
func (a *App) createNewNote() {
	newNote := &domain.Note{
		Title:     "",
		Content:   "",
		Version:   1,
		Metadata:  make(map[string]interface{}),
		CreatedAt: atom.Now(),
		UpdatedAt: atom.Now(),
	}

	a.mainUI.StartSaving()
	defer a.mainUI.EndSaving()
	a.currentNote = newNote
	a.hasUnsavedChanges = true // Mark as unsaved since it's not in database yet
	a.mainUI.DisplayNote(newNote)
	a.mainUI.MarkAsUnsaved()
}

// onDeleteNote is called when user wants to delete the current note
func (a *App) onDeleteNote() {
	if a.currentNote == nil {
		dialog.ShowInformation("No Note Selected", "Please select a note to delete", a.window)
		return
	}

	if a.currentNote.ID == "" {
		dialog.ShowInformation("Unsaved Note", "This note has not been saved yet and cannot be deleted", a.window)
		return
	}

	dialog.ShowConfirm("Delete Note",
		"Are you sure you want to delete this note?",
		func(confirmed bool) {
			if confirmed {
				rail := flow.EmptyRail()
				err := a.noteService.DeleteNote(rail, a.currentNote.ID)
				if err != nil {
					dialog.ShowError(err, a.window)
					return
				}

				a.currentNote = nil
				a.hasUnsavedChanges = false
				a.mainUI.RefreshNoteList()

				lastNote, err := a.noteService.GetLastModifiedNote(rail)
				if err != nil {
					a.mainUI.StartSaving()
					a.mainUI.ShowEmptyState()
					a.mainUI.EndSaving()
					a.mainUI.MarkAsSaved()
				} else {
					a.mainUI.StartSaving()
					a.currentNote = lastNote
					a.mainUI.DisplayNote(lastNote)
					a.mainUI.EndSaving()
					a.mainUI.MarkAsSaved()
				}
			}
		},
		a.window,
	)
}

// onImportNote is called when user wants to import notes
func (a *App) onImportNote() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		path := reader.URI().Path()

		dialog.ShowConfirm("Duplicate Notes",
			"If notes with the same ID exist, do you want to overwrite them?",
			func(overwrite bool) {
				rail := flow.EmptyRail()
				notes, err := a.importExportService.ImportNotesFromFile(rail, path, func(note *domain.Note) bool {
					return overwrite
				})
				if err != nil {
					dialog.ShowError(err, a.window)
					return
				}

				a.mainUI.RefreshNoteList()
				dialog.ShowInformation("Import Successful", fmt.Sprintf("Successfully imported %d notes", len(notes)), a.window)
			},
			a.window,
		)
	}, a.window)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	fd.Show()
}

// onExportNote is called when user wants to export all notes
func (a *App) onExportNote() {
	rail := flow.EmptyRail()
	notes, err := a.noteService.ListNotes(rail)
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	if len(notes) == 0 {
		dialog.ShowInformation("No Notes", "There are no notes to export", a.window)
		return
	}

	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		// Create export structure with all notes
		type ExportData struct {
			Version int               `json:"version"`
			Notes   []domain.NoteJSON `json:"notes"`
			Count   int               `json:"count"`
		}

		exportData := ExportData{
			Version: 1,
			Notes:   make([]domain.NoteJSON, len(notes)),
			Count:   len(notes),
		}

		for i, note := range notes {
			exportData.Notes[i] = note.ToJSON()
		}

		data, err := json.MarshalIndent(exportData, "", "  ")
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		_, err = writer.Write(data)
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		dialog.ShowInformation("Export Successful", fmt.Sprintf("Successfully exported %d notes", len(notes)), a.window)
	}, a.window)

	fd.SetFileName("nota_export.json")
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	fd.Show()
}

// onSearch is called when user searches for notes
func (a *App) onSearch(query string) {
	if query == "" {
		a.mainUI.RefreshNoteList()
		return
	}

	rail := flow.EmptyRail()
	notes, err := a.noteService.SearchNotes(rail, query)
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	a.mainUI.DisplaySearchResults(notes)
}

// onPinNote is called when user toggles pin mode
func (a *App) onPinNote(pin bool) {
	if a.currentNote == nil {
		return
	}

	a.mainUI.SetPinned(pin)
	a.mainUI.ToggleMinimizedMode(pin)

	// Resize window for minimized mode
	if pin {
		a.window.Resize(fyne.NewSize(400, 300))
	} else {
		a.window.Resize(fyne.NewSize(1000, 800))
	}
}

// GetDatabaseLocation returns the database location
func (a *App) GetDatabaseLocation() string {
	return infrastructure.GetDatabaseLocation()
}

// ListNotes returns all notes
func (a *App) ListNotes() ([]*domain.Note, error) {
	rail := flow.EmptyRail()
	return a.noteService.ListNotes(rail)
}

// OnNoteSelected implements NoteSelectionHandler interface
func (a *App) OnNoteSelected(note *domain.Note) {
	a.onNoteSelected(note)
}

// OnContentChanged implements NoteEditHandler interface
func (a *App) OnContentChanged() {
	a.onContentChanged()
}

// OnSave implements NoteEditHandler interface
func (a *App) OnSave() {
	a.saveCurrentNote()
}

// OnCreateNote implements AppActionsHandler interface
func (a *App) OnCreateNote() {
	a.onCreateNote()
}

// OnDeleteNote implements AppActionsHandler interface
func (a *App) OnDeleteNote() {
	a.onDeleteNote()
}

// OnImportNote implements AppActionsHandler interface
func (a *App) OnImportNote() {
	a.onImportNote()
}

// OnExportNote implements AppActionsHandler interface
func (a *App) OnExportNote() {
	a.onExportNote()
}

// OnSearch implements SearchHandler interface
func (a *App) OnSearch(query string) {
	a.onSearch(query)
}

// OnPinNote implements PinHandler interface
func (a *App) OnPinNote(pin bool) {
	a.onPinNote(pin)
}
