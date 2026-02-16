package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/curtisnewbie/nota/internal/domain"
	"github.com/curtisnewbie/nota/internal/i18n"
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
	// Pagination fields
	currentOffset int
	pageSize      int
	currentQuery  string
	hasMore       bool
	loading       bool
	loadMoreBtn   *widget.Button
	// Auto-loading fields
	scrollContainer   *container.Scroll
	checkScrollTicker *time.Ticker
	checkScrollDone   chan bool
}

// NewNoteList creates a new note list
func NewNoteList(selectionHandler NoteSelectionHandler, searchHandler SearchHandler) *NoteList {
	return &NoteList{
		selectionHandler: selectionHandler,
		searchHandler:    searchHandler,
		currentOffset:    0,
		pageSize:         30,
		currentQuery:     "",
		hasMore:          true,
		loading:          false,
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
	t := i18n.T()

	n.searchEntry = widget.NewEntry()
	n.searchEntry.SetPlaceHolder(t.Editor.PlaceholderSearch)
	n.searchEntry.OnChanged = func(query string) {
		if n.searchHandler != nil {
			n.searchHandler.OnSearch(query)
		}
	}

	// Create toolbar with search only
	toolbar := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		n.searchEntry,
	)

	// Create widget.List for displaying notes
	n.noteList = widget.NewList(
		func() int { return len(n.notes) },
		func() fyne.CanvasObject {
			// Create a container for each note item
			titleLabel := widget.NewLabel("")
			dateLabel := widget.NewLabel("")
			dateLabel.TextStyle = fyne.TextStyle{Italic: true}
			return container.NewVBox(titleLabel, dateLabel)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= 0 && id < len(n.notes) {
				note := n.notes[id]
				container := obj.(*fyne.Container)
				titleLabel := container.Objects[0].(*widget.Label)
				dateLabel := container.Objects[1].(*widget.Label)
				titleLabel.SetText(note.Title)
				dateLabel.SetText(note.UpdatedAt.Format("2006/01/02"))
			}
		},
	)

	n.noteList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(n.notes) && n.selectionHandler != nil {
			n.selectionHandler.OnNoteSelected(n.notes[id])
		}
	}

	// Create Load More button (hidden - used for manual fallback if needed)
	n.loadMoreBtn = widget.NewButton("Load More", func() {
		n.loadMoreNotes()
	})
	n.loadMoreBtn.Hide()

	// Use Scroll container for the list
	n.scrollContainer = container.NewScroll(n.noteList)
	n.scrollContainer.SetMinSize(fyne.NewSize(150, 0)) // 0 height means fill available space

	// Don't start scroll checking - use Load More button instead
	// n.startScrollChecking()

	// Use Border layout to put Load More button at bottom
	n.rightPanel = container.NewBorder(
		toolbar,
		n.loadMoreBtn, // Load More button at bottom
		nil,
		nil,
		n.scrollContainer, // Fill all available space
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
	deleteItem := fyne.NewMenuItem(i18n.T().Menu.Delete, func() {
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

// DisplayNotes displays the list of notes (for initial load or refresh)
func (n *NoteList) DisplayNotes(notes []*domain.Note) {
	// This method is called by the app with all notes or initial batch
	// For pagination, we use LoadNotes or LoadMoreNotes instead
	if notes == nil {
		n.notes = []*domain.Note{}
	} else {
		n.notes = notes
		n.hasMore = len(notes) == n.pageSize
	}
	n.noteList.Refresh()
}

// LoadNotes loads notes from scratch with pagination
func (n *NoteList) LoadNotes(notes []*domain.Note) {
	n.currentOffset = len(notes)
	n.notes = notes
	n.hasMore = len(notes) == n.pageSize
	n.loading = false
	n.noteList.Refresh()

	// Show/hide Load More button
	if n.hasMore && len(notes) > 0 {
		n.loadMoreBtn.Show()
	} else {
		n.loadMoreBtn.Hide()
	}
}

// AppendNotes appends notes to the existing list (for pagination)
func (n *NoteList) AppendNotes(notes []*domain.Note) {
	if len(notes) > 0 {
		n.notes = append(n.notes, notes...)
		n.currentOffset = len(n.notes) // Use total count, not just the new notes
		n.hasMore = len(notes) == n.pageSize
		n.noteList.Refresh()
	} else {
		n.hasMore = false
	}
	n.loading = false

	// Show/hide Load More button
	if n.hasMore {
		n.loadMoreBtn.Show()
	} else {
		n.loadMoreBtn.Hide()
	}
}

// loadMoreNotes loads the next page of notes
func (n *NoteList) loadMoreNotes() {
	if n.loading || !n.hasMore {
		return
	}

	n.loading = true
	if n.searchHandler != nil {
		n.searchHandler.OnSearch(n.currentQuery)
	}
}

// SetLoading sets the loading state
func (n *NoteList) SetLoading(loading bool) {
	n.loading = loading
}

// IsLoading returns whether currently loading
func (n *NoteList) IsLoading() bool {
	return n.loading
}

// GetOffset returns the current offset for pagination
func (n *NoteList) GetOffset() int {
	return n.currentOffset
}

// GetPageSize returns the page size
func (n *NoteList) GetPageSize() int {
	return n.pageSize
}

// GetCurrentQuery returns the current search query
func (n *NoteList) GetCurrentQuery() string {
	return n.currentQuery
}

// SetCurrentQuery sets the current search query
func (n *NoteList) SetCurrentQuery(query string) {
	n.currentQuery = query
}

// HasMore returns whether there are more notes to load
func (n *NoteList) HasMore() bool {
	return n.hasMore
}

// SetHasMore sets whether there are more notes to load
func (n *NoteList) SetHasMore(hasMore bool) {
	n.hasMore = hasMore
}

// startScrollChecking starts the scroll position checking goroutine
func (n *NoteList) startScrollChecking() {
	if n.checkScrollTicker != nil {
		// Already running
		return
	}

	n.checkScrollTicker = time.NewTicker(200 * time.Millisecond)
	n.checkScrollDone = make(chan bool)

	go func() {
		for {
			select {
			case <-n.checkScrollTicker.C:
				// Check if we're on the main thread, if not we need to handle this differently
				// For now, just call directly - the checkScrollPosition method itself should be safe
				n.checkScrollPosition()
			case <-n.checkScrollDone:
				return
			}
		}
	}()
}

// checkScrollPosition checks if scrolled near bottom and triggers load more
func (n *NoteList) checkScrollPosition() {
	if !n.hasMore || n.loading || n.scrollContainer == nil {
		return
	}

	// Calculate scroll position
	offset := n.scrollContainer.Offset.Y
	contentSize := n.scrollContainer.Content.Size().Height
	size := n.scrollContainer.Size().Height

	// Load more when scrolled to within 100px of bottom
	threshold := float32(100)
	if offset >= float32(contentSize-size)-threshold {
		n.loadMoreNotes()
	}
}

// StopScrollChecking stops the scroll position checking goroutine
func (n *NoteList) StopScrollChecking() {
	if n.checkScrollTicker != nil {
		n.checkScrollTicker.Stop()
		close(n.checkScrollDone)
		n.checkScrollTicker = nil
		n.checkScrollDone = nil
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
