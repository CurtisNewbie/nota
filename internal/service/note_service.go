package service

import (
	"errors"
	"fmt"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/nota/internal/domain"
	"github.com/curtisnewbie/nota/internal/repository"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrEmptyTitle   = errors.New("title cannot be empty")
)

// NoteService defines the interface for note business operations
type NoteService interface {
	CreateNote(rail flow.Rail, note *domain.Note) error
	UpdateNote(rail flow.Rail, note *domain.Note) error
	DeleteNote(rail flow.Rail, id string) error
	GetNote(rail flow.Rail, id string) (*domain.Note, error)
	ListNotes(rail flow.Rail) ([]*domain.Note, error)
	SearchNotes(rail flow.Rail, query string) ([]*domain.Note, error)
	GetLastModifiedNote(rail flow.Rail) (*domain.Note, error)
}

// NoteServiceImpl implements NoteService
type NoteServiceImpl struct {
	noteRepo repository.NoteRepository
}

// NewNoteService creates a new note service
func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &NoteServiceImpl{noteRepo: noteRepo}
}

// CreateNote creates a new note
func (s *NoteServiceImpl) CreateNote(rail flow.Rail, note *domain.Note) error {
	rail.Infof("Creating new note: %s", note.Title)
	if note.Title == "" {
		rail.Warnf("Attempted to create note with empty title")
		return ErrEmptyTitle
	}
	err := s.noteRepo.Save(rail, note)
	if err != nil {
		rail.Errorf("Failed to create note: %v", err)
	} else {
		rail.Infof("Successfully created note: %s", note.ID)
	}
	return err
}

// UpdateNote updates an existing note
func (s *NoteServiceImpl) UpdateNote(rail flow.Rail, note *domain.Note) error {
	rail.Infof("Updating note: %s", note.ID)
	if note.ID == "" {
		rail.Warnf("Attempted to update note with empty ID")
		return fmt.Errorf("note ID cannot be empty")
	}
	if note.Title == "" {
		rail.Warnf("Attempted to update note with empty title")
		return ErrEmptyTitle
	}

	_, err := s.noteRepo.FindByID(rail, note.ID)
	if err != nil {
		rail.Warnf("Note not found for update: %s", note.ID)
		return ErrNoteNotFound
	}

	err = s.noteRepo.Save(rail, note)
	if err != nil {
		rail.Errorf("Failed to update note: %v", err)
	} else {
		rail.Infof("Successfully updated note: %s", note.ID)
	}
	return err
}

// DeleteNote soft-deletes a note by ID
func (s *NoteServiceImpl) DeleteNote(rail flow.Rail, id string) error {
	rail.Infof("Deleting note: %s", id)
	if id == "" {
		rail.Warnf("Attempted to delete note with empty ID")
		return fmt.Errorf("note ID cannot be empty")
	}

	_, err := s.noteRepo.FindByID(rail, id)
	if err != nil {
		rail.Warnf("Note not found for deletion: %s", id)
		return ErrNoteNotFound
	}

	return s.noteRepo.Delete(rail, id)
}

// GetNote retrieves a single note by ID
func (s *NoteServiceImpl) GetNote(rail flow.Rail, id string) (*domain.Note, error) {
	rail.Debugf("Getting note: %s", id)
	if id == "" {
		rail.Warnf("Attempted to get note with empty ID")
		return nil, fmt.Errorf("note ID cannot be empty")
	}

	note, err := s.noteRepo.FindByID(rail, id)
	if err != nil {
		rail.Warnf("Note not found: %s", id)
		return nil, ErrNoteNotFound
	}

	return note, nil
}

// ListNotes retrieves all notes (excludes soft-deleted), sorted by updated_at DESC
func (s *NoteServiceImpl) ListNotes(rail flow.Rail) ([]*domain.Note, error) {
	rail.Debugf("Listing all notes")
	return s.noteRepo.FindAllSorted(rail)
}

// SearchNotes searches notes by title and content using FTS
func (s *NoteServiceImpl) SearchNotes(rail flow.Rail, query string) ([]*domain.Note, error) {
	rail.Debugf("Searching notes with query: %s", query)
	return s.noteRepo.Search(rail, query)
}

// GetLastModifiedNote retrieves the most recently modified note
func (s *NoteServiceImpl) GetLastModifiedNote(rail flow.Rail) (*domain.Note, error) {
	rail.Debugf("Getting last modified note")
	notes, err := s.noteRepo.FindAllSorted(rail)
	if err != nil {
		rail.Errorf("Failed to get last modified note: %v", err)
		return nil, err
	}

	if len(notes) == 0 {
		rail.Warnf("No notes found")
		return nil, ErrNoteNotFound
	}

	rail.Infof("Last modified note: %s", notes[0].ID)
	return notes[0], nil
}
