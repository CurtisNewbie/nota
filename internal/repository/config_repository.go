package repository

import (
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/nota/internal/domain"
	"gorm.io/gorm"
)

// ConfigRepository defines the interface for config data operations
type ConfigRepository interface {
	Save(rail flow.Rail, config *domain.Config) error
	FindByName(rail flow.Rail, name string) (*domain.Config, error)
	FindAll(rail flow.Rail) ([]*domain.Config, error)
	Delete(rail flow.Rail, name string) error
}

// SQLiteConfigRepository implements ConfigRepository for SQLite
type SQLiteConfigRepository struct {
	db *gorm.DB
}

// NewSQLiteConfigRepository creates a new SQLite config repository
func NewSQLiteConfigRepository(db *gorm.DB) ConfigRepository {
	return &SQLiteConfigRepository{db: db}
}

// Save saves a config (create or update)
func (r *SQLiteConfigRepository) Save(rail flow.Rail, config *domain.Config) error {
	rail.Debugf("Saving config: %s", config.Name)
	q := dbquery.NewQuery(rail, r.db).Table("config")
	_, err := q.Create(config)
	if err != nil {
		rail.Errorf("Failed to save config %s: %v", config.Name, err)
	} else {
		rail.Infof("Successfully saved config: %s", config.Name)
	}
	return err
}

// FindByName finds a config by name
func (r *SQLiteConfigRepository) FindByName(rail flow.Rail, name string) (*domain.Config, error) {
	rail.Debugf("Finding config by name: %s", name)
	var config domain.Config
	q := dbquery.NewQuery(rail, r.db).Table("config").Where("name = ?", name)
	_, err := q.Scan(&config)
	if err != nil {
		rail.Warnf("Config not found: %s", name)
		return nil, err
	}
	return &config, nil
}

// FindAll finds all configs
func (r *SQLiteConfigRepository) FindAll(rail flow.Rail) ([]*domain.Config, error) {
	rail.Debugf("Finding all configs")
	var configs []*domain.Config
	q := dbquery.NewQuery(rail, r.db).Table("config")
	_, err := q.Scan(&configs)
	rail.Debugf("Found %d configs", len(configs))
	return configs, err
}

// Delete deletes a config by name
func (r *SQLiteConfigRepository) Delete(rail flow.Rail, name string) error {
	rail.Infof("Deleting config: %s", name)
	q := dbquery.NewQuery(rail, r.db).Table("config").Where("name = ?", name)
	_, err := q.Delete()
	if err != nil {
		rail.Errorf("Failed to delete config %s: %v", name, err)
	} else {
		rail.Infof("Successfully deleted config: %s", name)
	}
	return err
}
