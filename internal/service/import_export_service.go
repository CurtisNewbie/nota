package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/nota/internal/domain"
	"github.com/curtisnewbie/nota/internal/repository"
)

// ImportExportService defines the interface for import/export operations
type ImportExportService interface {
	ExportNote(note *domain.Note, path string) error
	ExportNotes(notes []*domain.Note, dir string) error
	ImportNote(path string, onDuplicate func(note *domain.Note) bool) (*domain.Note, error)
	ImportNotes(dir string, onDuplicate func(note *domain.Note) bool) ([]*domain.Note, error)
}

// ImportExportServiceImpl implements ImportExportService
type ImportExportServiceImpl struct {
	noteRepo repository.NoteRepository
}

// NewImportExportService creates a new import/export service
func NewImportExportService(noteRepo repository.NoteRepository) ImportExportService {
	return &ImportExportServiceImpl{noteRepo: noteRepo}
}

// ExportNote exports a single note to a JSON file
func (s *ImportExportServiceImpl) ExportNote(note *domain.Note, path string) error {
	rail := flow.NewRail(context.Background())
	
	if note == nil {
		return fmt.Errorf("note cannot be nil")
	}
	
	noteJSON := note.ToJSON()
	
	data, err := json.MarshalIndent(noteJSON, "", "  ")
	if err != nil {
		rail.Errorf("Failed to marshal note JSON: %v", err)
		return fmt.Errorf("failed to marshal note: %w", err)
	}
	
	if path == "" {
		timestamp := time.Now().Format("20060102_150405")
		path = fmt.Sprintf("%s_nota_exported.json", timestamp)
	}
	
	if filepath.Ext(path) != ".json" {
		path += ".json"
	}
	
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		rail.Errorf("Failed to write note to file: %v", err)
		return fmt.Errorf("failed to write note file: %w", err)
	}
	
	rail.Infof("Successfully exported note to: %s", path)
	return nil
}

// ExportNotes exports multiple notes to a directory (batch operation)
func (s *ImportExportServiceImpl) ExportNotes(notes []*domain.Note, dir string) error {
	rail := flow.NewRail(context.Background())
	
	if len(notes) == 0 {
		return fmt.Errorf("no notes to export")
	}
	
	if dir == "" {
		return fmt.Errorf("directory path cannot be empty")
	}
	
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		rail.Errorf("Failed to create export directory: %v", err)
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	successCount := 0
	for _, note := range notes {
		timestamp := note.UpdatedAt.Format("20060102_150405")
		filename := fmt.Sprintf("%s_%s_nota_exported.json", timestamp, note.ID[:8])
		path := filepath.Join(dir, filename)
		
		err := s.ExportNote(note, path)
		if err != nil {
			rail.Warnf("Failed to export note %s: %v", note.ID, err)
			continue
		}
		successCount++
	}
	
	rail.Infof("Successfully exported %d/%d notes to: %s", successCount, len(notes), dir)
	
	if successCount == 0 {
		return fmt.Errorf("failed to export any notes")
	}
	
	return nil
}

// ImportNote imports a single note from a JSON file
func (s *ImportExportServiceImpl) ImportNote(path string, onDuplicate func(note *domain.Note) bool) (*domain.Note, error) {
	rail := flow.NewRail(context.Background())
	
	if path == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		rail.Errorf("Failed to read note file: %v", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	var noteJSON domain.NoteJSON
	err = json.Unmarshal(data, &noteJSON)
	if err != nil {
		rail.Errorf("Failed to unmarshal note JSON: %v", err)
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	if noteJSON.Version != 1 {
		rail.Warnf("Unsupported note version: %d", noteJSON.Version)
		return nil, fmt.Errorf("unsupported note version: %d", noteJSON.Version)
	}
	
	note, err := domain.FromJSON(noteJSON)
	if err != nil {
		rail.Errorf("Failed to convert JSON to note: %v", err)
		return nil, fmt.Errorf("failed to convert note: %w", err)
	}
	
	existing, err := s.noteRepo.FindByID(note.ID)
	if err == nil {
		if onDuplicate != nil && onDuplicate(existing) {
			err = s.noteRepo.Save(note)
			if err != nil {
				rail.Errorf("Failed to overwrite note: %v", err)
				return nil, fmt.Errorf("failed to overwrite note: %w", err)
			}
			rail.Infof("Successfully overwrote note: %s", note.ID)
		} else {
			rail.Infof("Skipped duplicate note: %s", note.ID)
			return existing, nil
		}
	} else {
		err = s.noteRepo.Save(note)
		if err != nil {
			rail.Errorf("Failed to save note: %v", err)
			return nil, fmt.Errorf("failed to save note: %w", err)
		}
		rail.Infof("Successfully imported note: %s", note.ID)
	}
	
	return note, nil
}

// ImportNotes imports multiple notes from a directory (batch operation)
func (s *ImportExportServiceImpl) ImportNotes(dir string, onDuplicate func(note *domain.Note) bool) ([]*domain.Note, error) {
	rail := flow.NewRail(context.Background())
	
	if dir == "" {
		return nil, fmt.Errorf("directory path cannot be empty")
	}
	
	files, err := os.ReadDir(dir)
	if err != nil {
		rail.Errorf("Failed to read directory: %v", err)
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	
	var importedNotes []*domain.Note
	successCount := 0
	
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		path := filepath.Join(dir, file.Name())
		note, err := s.ImportNote(path, onDuplicate)
		if err != nil {
			rail.Warnf("Failed to import note from %s: %v", file.Name(), err)
			continue
		}
		importedNotes = append(importedNotes, note)
		successCount++
	}
	
	rail.Infof("Successfully imported %d notes from: %s", successCount, dir)
	
	if len(importedNotes) == 0 {
		return nil, fmt.Errorf("no valid notes found to import")
	}
	
	return importedNotes, nil
}