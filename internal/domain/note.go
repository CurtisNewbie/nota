package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/idutil"
	"gorm.io/gorm"
)

// Note represents a note in the system
type Note struct {
	ID        string                 `gorm:"primaryKey" json:"id"`
	Title     string                 `gorm:"not null" json:"title"`
	Content   string                 `gorm:"type:text" json:"content"`
	Version   int                    `gorm:"not null;default:1" json:"version"`
	CreatedAt atom.Time              `gorm:"not null" json:"created_at"`
	UpdatedAt atom.Time              `gorm:"not null" json:"updated_at"`
	DeletedAt *atom.Time             `gorm:"index" json:"deleted_at,omitempty"`
	Metadata  map[string]interface{} `gorm:"type:text;serializer:json" json:"metadata"`
}

// TableName specifies the table name for GORM
func (Note) TableName() string {
	return "note"
}

// BeforeCreate GORM hook to generate ID and set timestamps
func (n *Note) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = idutil.Id("note")
	}
	now := atom.WrapTime(time.Now())
	n.CreatedAt = now
	n.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM hook to update timestamp
func (n *Note) BeforeUpdate(tx *gorm.DB) error {
	n.UpdatedAt = atom.WrapTime(time.Now())
	return nil
}

// NoteJSON represents the JSON format for import/export
type NoteJSON struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Version   int                    `json:"version"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	DeletedAt *string                `json:"deleted_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ToJSON converts Note to NoteJSON for export
func (n *Note) ToJSON() NoteJSON {
	result := NoteJSON{
		ID:        n.ID,
		Title:     n.Title,
		Content:   n.Content,
		Version:   n.Version,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
		UpdatedAt: n.UpdatedAt.Format(time.RFC3339),
		Metadata:  n.Metadata,
	}
	if n.DeletedAt != nil {
		deletedAt := n.DeletedAt.Format(time.RFC3339)
		result.DeletedAt = &deletedAt
	}
	return result
}

// FromJSON creates Note from NoteJSON for import
func FromJSON(json NoteJSON) (*Note, error) {
	note := &Note{
		ID:       json.ID,
		Title:    json.Title,
		Content:  json.Content,
		Version:  json.Version,
		Metadata: json.Metadata,
	}

	if json.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, json.CreatedAt); err == nil {
			note.CreatedAt = atom.WrapTime(t)
		} else {
			note.CreatedAt = atom.WrapTime(time.Now())
		}
	} else {
		note.CreatedAt = atom.WrapTime(time.Now())
	}

	if json.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, json.UpdatedAt); err == nil {
			note.UpdatedAt = atom.WrapTime(t)
		} else {
			note.UpdatedAt = atom.WrapTime(time.Now())
		}
	} else {
		note.UpdatedAt = atom.WrapTime(time.Now())
	}

	if json.DeletedAt != nil && *json.DeletedAt != "" {
		if t, err := time.Parse(time.RFC3339, *json.DeletedAt); err == nil {
			deletedAt := atom.WrapTime(t)
			note.DeletedAt = &deletedAt
		}
	}

	return note, nil
}

// JSONMap is a custom type for storing JSON in SQLite
type JSONMap map[string]interface{}

// Value implements driver.Valuer for JSONMap
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner for JSONMap
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}
