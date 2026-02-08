package repository

import (
	"github.com/curtisnewbie/nota/internal/domain"
	"gorm.io/gorm"
)

// ConfigRepository defines the interface for config data operations
type ConfigRepository interface {
	Save(config *domain.Config) error
	FindByName(name string) (*domain.Config, error)
	FindAll() ([]*domain.Config, error)
	Delete(name string) error
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
func (r *SQLiteConfigRepository) Save(config *domain.Config) error {
	return r.db.Save(config).Error
}

// FindByName finds a config by name
func (r *SQLiteConfigRepository) FindByName(name string) (*domain.Config, error) {
	var config domain.Config
	err := r.db.Where("name = ?", name).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// FindAll finds all configs
func (r *SQLiteConfigRepository) FindAll() ([]*domain.Config, error) {
	var configs []*domain.Config
	err := r.db.Find(&configs).Error
	return configs, err
}

// Delete deletes a config by name
func (r *SQLiteConfigRepository) Delete(name string) error {
	return r.db.Where("name = ?", name).Delete(&domain.Config{}).Error
}
