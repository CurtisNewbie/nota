package infrastructure

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/nota/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	defaultDatabasePath = "$HOME/nota/data/nota.sqlite"
)

// InitializeDatabase initializes the SQLite database with schema migration
func InitializeDatabase() (*gorm.DB, error) {
	rail := flow.NewRail(context.Background())
	
	dbPath := getDatabasePath()
	
	rail.Infof("Initializing database at: %s", dbPath)
	
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		rail.Errorf("Failed to open database: %v", err)
		return nil, err
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		rail.Errorf("Failed to get database instance: %v", err)
		return nil, err
	}
	
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	
	err = db.Exec("PRAGMA foreign_keys = ON").Error
	if err != nil {
		rail.Errorf("Failed to enable foreign keys: %v", err)
		return nil, err
	}
	
	err = db.AutoMigrate(&domain.Note{}, &domain.Config{})
	if err != nil {
		rail.Errorf("Failed to migrate database schema: %v", err)
		return nil, err
	}
	
	err = createFTS5Table(db)
	if err != nil {
		rail.Errorf("Failed to create FTS5 table: %v", err)
		return nil, err
	}
	
	dbquery.ImplGetPrimaryDBFunc(func() *gorm.DB {
		return db
	})
	
	rail.Infof("Database initialized successfully")
	
	return db, nil
}

// getDatabasePath returns the database path from config or default
func getDatabasePath() string {
	return os.ExpandEnv(defaultDatabasePath)
}

// createFTS5Table creates the full-text search virtual table
func createFTS5Table(db *gorm.DB) error {
	rail := flow.NewRail(context.Background())
	
	// Try to create FTS5 table, but if it fails (module not available), continue without it
	err := db.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS note_fts USING fts5(title, content);
	`).Error
	if err != nil {
		rail.Warnf("FTS5 module not available, falling back to LIKE-based search: %v", err)
		return nil // Continue without FTS5
	}
	
	err = db.Exec(`
		CREATE TRIGGER IF NOT EXISTS note_ai AFTER INSERT ON note BEGIN
			INSERT INTO note_fts(rowid, title, content) VALUES (new.rowid, new.title, new.content);
		END;
	`).Error
	if err != nil {
		rail.Warnf("Failed to create FTS5 insert trigger: %v", err)
	}
	
	err = db.Exec(`
		CREATE TRIGGER IF NOT EXISTS note_ad AFTER DELETE ON note BEGIN
			INSERT INTO note_fts(note_fts, rowid, title, content) VALUES('delete', old.rowid, old.title, old.content);
		END;
	`).Error
	if err != nil {
		rail.Warnf("Failed to create FTS5 delete trigger: %v", err)
	}
	
	err = db.Exec(`
		CREATE TRIGGER IF NOT EXISTS note_au AFTER UPDATE ON note BEGIN
			INSERT INTO note_fts(note_fts, rowid, title, content) VALUES('delete', old.rowid, old.title, old.content);
			INSERT INTO note_fts(rowid, title, content) VALUES (new.rowid, new.title, new.content);
		END;
	`).Error
	if err != nil {
		rail.Warnf("Failed to create FTS5 update trigger: %v", err)
	}
	
	rail.Infof("FTS5 table and triggers created successfully")
	return nil
}

// EnsureDatabaseDir ensures the database directory exists
func EnsureDatabaseDir() error {
	dbPath := getDatabasePath()
	dir := filepath.Dir(dbPath)
	
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
	}
	
	return nil
}

// GetDatabaseLocation returns the current database location
func GetDatabaseLocation() string {
	return getDatabasePath()
}