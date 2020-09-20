package std

import (
	"testing"

	"github.com/beatlabs/patron/log"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	logger := Create(log.InfoLevel)(map[string]interface{}{"name": "john"})
	assert.NotNil(t, logger)
}
