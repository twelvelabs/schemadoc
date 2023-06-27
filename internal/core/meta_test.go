package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMeta(t *testing.T) {
	meta := NewMeta("1.2.3", "9b11774", "2023-02-19T00:57:51-06:00")
	assert.Equal(t, "9b11774", meta.BuildCommit)
	assert.Equal(t, "2023-02-19T00:57:51-06:00", meta.BuildTime.Format(time.RFC3339))
	assert.Equal(t, "1.2.3", meta.Version)
}
