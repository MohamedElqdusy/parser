package process

import (
	"github.com/stretchr/testify/assert"
	"parser/logger"
	"parser/resources"
	"testing"
)

func TestLocalFileProcessor(t *testing.T) {
	log := logger.NewLogger()
	logger.InitLogger(log)
	defer logger.Sync()

	// creates a new processor for a local file.
	p, err := NewFileProcessor("example")
	assert.NoError(t, err)

	// starts processing
	err = p.Process(&resources.Customer{})
	assert.NoError(t, err)

	// check the parsing summary stats
	assert.Equal(t, int64(11), p.Summary.Total)
	assert.Equal(t, int64(5), p.Summary.Skipped)
	assert.Equal(t, int64(5), p.Summary.Valid)
}
