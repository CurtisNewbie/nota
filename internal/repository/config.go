package repository

import (
	"time"

	"github.com/curtisnewbie/miso/middleware/sqlite"
	"gorm.io/gorm/logger"
)

func init() {
	sqlite.UpdateLoggerConfig(logger.Config{SlowThreshold: 200 * time.Millisecond, LogLevel: logger.Info, Colorful: false})
}
