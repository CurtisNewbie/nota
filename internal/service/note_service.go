package service

import (
	"errors"
	"fmt"

	"github.com/curtisnewbie/nota/internal/domain"
	"github.com/curtisnewbie/nota/internal/repository"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrEmptyTitle   = errors.New("title cannot be empty")
)

// NoteService defines the interface for note business operations
type NoteService interface {
	CreateNote(note *domain.Note) error
	UpdateNote(note *domain.Note) error
	DeleteNote(id string) error
	GetNote(id string) (*domain.Note, error)
	ListNotes() ([]*domain.Note, error)
	SearchNotes(query string) ([]*domain.Note, error)
	GetLastModifiedNote() (*domain.Note, error)
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
func (s *NoteServiceImpl) CreateNote(note *domain.Note) error {
	if note.Title == "" {
		return ErrEmptyTitle
	}
	return s.noteRepo.Save(note)
}

// UpdateNote updates an existing note
func (s *NoteServiceImpl) UpdateNote(note *domain.Note) error {
	if note.ID == "" {
		return fmt.Errorf("note ID cannot be empty")
	}
	if note.Title == "" {
		return ErrEmptyTitle
	}
	
	_, err := s.noteRepo.FindByID(note.ID)
	if err != nil {
		return ErrNoteNotFound
	}
	
	return s.noteRepo.Save(note)
}

// DeleteNote soft-deletes a note by ID
func (s *NoteServiceImpl) DeleteNote(id string) error {
	if id == "" {
		return fmt.Errorf("note ID cannot be empty")
	}
	
	_, err := s.noteRepo.FindByID(id)
	if err != nil {
		return ErrNoteNotFound
	}
	
	return s.noteRepo.Delete(id)
}

// GetNote retrieves a single note by ID
func (s *NoteServiceImpl) GetNote(id string) (*domain.Note, error) {
	if id == "" {
		return nil, fmt.Errorf("note ID cannot be empty")
	}
	
	note, err := s.noteRepo.FindByID(id)
	if err != nil {
		return nil, ErrNoteNotFound
	}
	
	return note, nil
}

// ListNotes retrieves all notes (excludes soft-deleted), sorted by updated_at DESC
func (s *NoteServiceImpl) ListNotes() ([]*domain.Note, error) {
	return s.noteRepo.FindAllSorted()
}

// SearchNotes searches notes by title and content using FTS
func (s *NoteServiceImpl) SearchNotes(query string) ([]*domain.Note, error) {
	return s.noteRepo.Search(query)
}

// GetLastModifiedNote retrieves the most recently modified note
func (s *NoteServiceImpl) GetLastModifiedNote() (*domain.Note, error) {
	notes, err := s.noteRepo.FindAllSorted()
	if err != nil {
		return nil, err
	}
	
	if len(notes) == 0 {
		return nil, ErrNoteNotFound
	}
	
	return notes[0], nil
}