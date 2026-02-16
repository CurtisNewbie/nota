package repository

import (
	"time"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/idutil"
	"github.com/curtisnewbie/nota/internal/domain"
	"gorm.io/gorm"
)

// NoteRepository defines the interface for note data operations
type NoteRepository interface {
	Save(rail flow.Rail, note *domain.Note) error
	FindByID(rail flow.Rail, id string) (*domain.Note, error)
	FindAll(rail flow.Rail) ([]*domain.Note, error)
	FindAllSorted(rail flow.Rail) ([]*domain.Note, error)
	FindAllSortedPaginated(rail flow.Rail, offset, limit int) ([]*domain.Note, error)
	Search(rail flow.Rail, query string) ([]*domain.Note, error)
	SearchPaginated(rail flow.Rail, query string, offset, limit int) ([]*domain.Note, error)
	Delete(rail flow.Rail, id string) error
	FindByTitle(rail flow.Rail, title string) (*domain.Note, error)
	FindLastModified(rail flow.Rail) (*domain.Note, error)
}

// SQLiteNoteRepository implements NoteRepository for SQLite
type SQLiteNoteRepository struct {
	db *gorm.DB
}

// NewSQLiteNoteRepository creates a new SQLite note repository
func NewSQLiteNoteRepository(db *gorm.DB) NoteRepository {
	return &SQLiteNoteRepository{db: db}
}

// Save saves or updates a note
func (r *SQLiteNoteRepository) Save(rail flow.Rail, note *domain.Note) error {
	if note.ID == "" {
		note.ID = idutil.Id("note")
		q := dbquery.NewQuery(rail, r.db).Table("note")
		_, err := q.Create(note)
		return err
	}
	// For updates, use Set to specify columns
	q := dbquery.NewQuery(rail, r.db).Table("note").Where("id = ?", note.ID).Set("title", note.Title).Set("content", note.Content).Set("updated_at", atom.Now())
	_, err := q.Update()
	return err
}

// FindByID finds a note by ID
func (r *SQLiteNoteRepository) FindByID(rail flow.Rail, id string) (*domain.Note, error) {
	rail.Debugf("Finding note by ID: %s", id)
	var note domain.Note
	q := dbquery.NewQuery(rail, r.db).Table("note").Where("id = ?", id)
	_, err := q.Scan(&note)
	if err != nil {
		rail.Warnf("Note not found: %s", id)
		return nil, err
	}
	return &note, nil
}

// FindAll finds all notes (excluding soft-deleted)
func (r *SQLiteNoteRepository) FindAll(rail flow.Rail) ([]*domain.Note, error) {
	rail.Debugf("Finding all notes")
	var notes []*domain.Note
	q := dbquery.NewQuery(rail, r.db).Table("note").Where("deleted_at IS NULL")
	_, err := q.Scan(&notes)
	rail.Debugf("Found %d notes", len(notes))
	return notes, err
}

// FindAllSorted finds all notes sorted by updated_at DESC (excluding soft-deleted)
func (r *SQLiteNoteRepository) FindAllSorted(rail flow.Rail) ([]*domain.Note, error) {
	rail.Debugf("Finding all notes sorted by updated_at")
	var notes []*domain.Note
	q := dbquery.NewQuery(rail, r.db).Table("note").Where("deleted_at IS NULL").Order("updated_at DESC")
	_, err := q.Scan(&notes)
	rail.Debugf("Found %d notes", len(notes))
	return notes, err
}

// FindAllSortedPaginated finds notes sorted by updated_at DESC with pagination (excluding soft-deleted)
func (r *SQLiteNoteRepository) FindAllSortedPaginated(rail flow.Rail, offset, limit int) ([]*domain.Note, error) {
	rail.Debugf("Finding notes sorted by updated_at (offset=%d, limit=%d)", offset, limit)
	var notes []*domain.Note
	q := dbquery.NewQuery(rail, r.db).Table("note").
		Where("deleted_at IS NULL").
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset)
	_, err := q.Scan(&notes)
	rail.Debugf("Found %d notes", len(notes))
	return notes, err
}

// Search searches notes by title and content using LIKE-based search
func (r *SQLiteNoteRepository) Search(rail flow.Rail, query string) ([]*domain.Note, error) {
	if query == "" {
		return r.FindAllSorted(rail)
	}

	rail.Debugf("Searching notes with query: %s", query)
	var notes []*domain.Note
	searchPattern := "%" + query + "%"
	q := dbquery.NewQuery(rail, r.db).Table("note").
		Where("deleted_at IS NULL AND (title LIKE ? OR content LIKE ?)", searchPattern, searchPattern).
		Order("updated_at DESC")
	_, err := q.Scan(&notes)
	rail.Debugf("Found %d notes matching query", len(notes))
	return notes, err
}

// SearchPaginated searches notes by title and content using LIKE-based search with pagination
func (r *SQLiteNoteRepository) SearchPaginated(rail flow.Rail, query string, offset, limit int) ([]*domain.Note, error) {
	if query == "" {
		return r.FindAllSortedPaginated(rail, offset, limit)
	}

	rail.Debugf("Searching notes with query: %s (offset=%d, limit=%d)", query, offset, limit)
	var notes []*domain.Note
	searchPattern := "%" + query + "%"
	q := dbquery.NewQuery(rail, r.db).Table("note").
		Where("deleted_at IS NULL AND (title LIKE ? OR content LIKE ?)", searchPattern, searchPattern).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset)
	_, err := q.Scan(&notes)
	rail.Debugf("Found %d notes matching query", len(notes))
	return notes, err
}

// Delete soft-deletes a note by setting deleted_at timestamp
func (r *SQLiteNoteRepository) Delete(rail flow.Rail, id string) error {
	rail.Infof("Deleting note: %s", id)
	now := atom.WrapTime(time.Now())
	q := dbquery.NewQuery(rail, r.db).Table("note").Where("id = ?", id).Set("deleted_at", now)
	_, err := q.Update()
	if err != nil {
		rail.Errorf("Failed to delete note %s: %v", id, err)
	} else {
		rail.Infof("Successfully deleted note: %s", id)
	}
	return err
}

// FindByTitle finds a note by title
func (r *SQLiteNoteRepository) FindByTitle(rail flow.Rail, title string) (*domain.Note, error) {
	rail.Debugf("Finding note by title: %s", title)
	var note domain.Note
	q := dbquery.NewQuery(rail, r.db).Table("note").Where("title = ? AND deleted_at IS NULL", title)
	_, err := q.Scan(&note)
	if err != nil {
		rail.Warnf("Note not found with title: %s", title)
		return nil, err
	}
	return &note, nil
}

// FindLastModified finds the most recently modified note (excluding soft-deleted)
func (r *SQLiteNoteRepository) FindLastModified(rail flow.Rail) (*domain.Note, error) {
	rail.Debugf("Finding last modified note")
	var note domain.Note
	q := dbquery.NewQuery(rail, r.db).Table("note").
		Where("deleted_at IS NULL").
		Order("updated_at DESC").
		Limit(1)
	_, err := q.Scan(&note)
	if err != nil {
		rail.Warnf("No notes found")
		return nil, err
	}
	rail.Debugf("Last modified note: %s", note.ID)
	return &note, nil
}
