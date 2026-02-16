package service

import (
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/nota/internal/domain"
	"github.com/curtisnewbie/nota/internal/i18n"
	"github.com/curtisnewbie/nota/internal/repository"
)

const (
	configKeyLanguage = "language"
)

// ConfigService defines the interface for config operations
type ConfigService interface {
	SaveLanguage(rail flow.Rail, lang i18n.Language) error
	GetLanguage(rail flow.Rail) (i18n.Language, error)
}

// ConfigServiceImpl implements ConfigService
type ConfigServiceImpl struct {
	configRepo repository.ConfigRepository
}

// NewConfigService creates a new config service
func NewConfigService(configRepo repository.ConfigRepository) ConfigService {
	return &ConfigServiceImpl{configRepo: configRepo}
}

// SaveLanguage saves the language preference
func (s *ConfigServiceImpl) SaveLanguage(rail flow.Rail, lang i18n.Language) error {
	rail.Infof("Saving language preference: %s", lang)

	config := &domain.Config{
		Name:  configKeyLanguage,
		Value: string(lang),
	}

	err := s.configRepo.Save(rail, config)
	if err != nil {
		rail.Errorf("Failed to save language preference: %v", err)
		return err
	}

	rail.Infof("Successfully saved language preference: %s", lang)
	return nil
}

// GetLanguage retrieves the language preference
func (s *ConfigServiceImpl) GetLanguage(rail flow.Rail) (i18n.Language, error) {
	rail.Debugf("Getting language preference")

	config, err := s.configRepo.FindByName(rail, configKeyLanguage)
	if err != nil {
		rail.Warnf("Failed to get language preference: %v, using default", err)
		return i18n.LanguageEnglish, nil
	}

	lang := i18n.Language(config.Value)
	if lang == "" {
		lang = i18n.LanguageEnglish
	}

	rail.Infof("Language preference: %s", lang)
	return lang, nil
}
