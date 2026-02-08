package domain

// Config represents a configuration key-value pair stored in database
type Config struct {
	Name  string `gorm:"primaryKey" json:"name"`
	Value string `json:"value"`
}

// TableName specifies the table name for GORM
func (Config) TableName() string {
	return "config"
}