package repository

import (
	"time"

	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/nota/internal/domain"
	"gorm.io/gorm"
)

// NoteRepository defines the interface for note data operations
type NoteRepository interface {
	Save(note *domain.Note) error
	FindByID(id string) (*domain.Note, error)
	FindAll() ([]*domain.Note, error)
	FindAllSorted() ([]*domain.Note, error)
	Search(query string) ([]*domain.Note, error)
	Delete(id string) error
	FindByTitle(title string) (*domain.Note, error)
}

// SQLiteNoteRepository implements NoteRepository for SQLite
type SQLiteNoteRepository struct {
	db *gorm.DB
}

// NewSQLiteNoteRepository creates a new SQLite note repository
func NewSQLiteNoteRepository(db *gorm.DB) NoteRepository {
	return &SQLiteNoteRepository{db: db}
}

// Save saves a note (create or update)
func (r *SQLiteNoteRepository) Save(note *domain.Note) error {
	return r.db.Save(note).Error
}

// FindByID finds a note by ID
func (r *SQLiteNoteRepository) FindByID(id string) (*domain.Note, error) {
	var note domain.Note
	err := r.db.Where("id = ?", id).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// FindAll finds all notes (excluding soft-deleted)
func (r *SQLiteNoteRepository) FindAll() ([]*domain.Note, error) {
	var notes []*domain.Note
	err := r.db.Where("deleted_at IS NULL").Find(&notes).Error
	return notes, err
}

// FindAllSorted finds all notes sorted by updated_at DESC (excluding soft-deleted)
func (r *SQLiteNoteRepository) FindAllSorted() ([]*domain.Note, error) {
	var notes []*domain.Note
	err := r.db.Where("deleted_at IS NULL").Order("updated_at DESC").Find(&notes).Error
	return notes, err
}

// Search searches notes by title and content using FTS5 or LIKE-based search
func (r *SQLiteNoteRepository) Search(query string) ([]*domain.Note, error) {
	if query == "" {
		return r.FindAllSorted()
	}
	
	var notes []*domain.Note
	
	// Try FTS5 search first
	err := r.db.Raw(`
		SELECT note.* FROM note
		INNER JOIN note_fts ON note.rowid = note_fts.rowid
		WHERE note_fts MATCH ? AND note.deleted_at IS NULL
		ORDER BY note.updated_at DESC
	`, query).Scan(&notes).Error
	
	// If FTS5 fails, fall back to LIKE-based search
	if err != nil {
		searchPattern := "%" + query + "%"
		err = r.db.Where("deleted_at IS NULL AND (title LIKE ? OR content LIKE ?)", searchPattern, searchPattern).
			Order("updated_at DESC").
			Find(&notes).Error
	}
	
	return notes, err
}

// Delete soft-deletes a note by setting deleted_at timestamp
func (r *SQLiteNoteRepository) Delete(id string) error {
	now := atom.WrapTime(time.Now())
	return r.db.Model(&domain.Note{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// FindByTitle finds a note by title
func (r *SQLiteNoteRepository) FindByTitle(title string) (*domain.Note, error) {
	var note domain.Note
	err := r.db.Where("title = ? AND deleted_at IS NULL", title).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}