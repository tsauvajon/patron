package std

import (
	"os"

	"github.com/beatlabs/patron/log"
)

// Create creates a zerolog factory with default settings.
func Create(lvl log.Level) log.FactoryFunc {
	return func(f map[string]interface{}) log.Logger {
		return NewLogger(os.Stderr, lvl, f)
	}
}
