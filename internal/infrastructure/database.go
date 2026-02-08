package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/sqlite"
	"github.com/curtisnewbie/nota/internal/domain"
	"gorm.io/gorm"
)

const (
	defaultDatabasePath = "$HOME/nota/data/nota.sqlite"
)

// InitializeDatabase initializes the SQLite database with schema migration
func InitializeDatabase() (*gorm.DB, error) {
	rail := flow.EmptyRail()

	dbPath := getDatabasePath()

	rail.Infof("Initializing database at: %s", dbPath)

	gormDB, err := sqlite.NewConn(dbPath, true)
	if err != nil {
		rail.Errorf("Failed to open database: %v", err)
		return nil, err
	}

	err = gormDB.Exec("PRAGMA foreign_keys = ON").Error
	if err != nil {
		rail.Errorf("Failed to enable foreign keys: %v", err)
		return nil, err
	}

	err = gormDB.AutoMigrate(&domain.Note{}, &domain.Config{})
	if err != nil {
		rail.Errorf("Failed to migrate database schema: %v", err)
		return nil, err
	}

	dbquery.ImplGetPrimaryDBFunc(func() *gorm.DB {
		return gormDB
	})

	rail.Infof("Database initialized successfully")

	return gormDB, nil
}

// getDatabasePath returns the database path from config or default
func getDatabasePath() string {
	return os.ExpandEnv(defaultDatabasePath)
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
