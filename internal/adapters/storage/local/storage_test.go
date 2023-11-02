package local

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_SessionStart(t *testing.T) {
	storage := NewStorage(NewConfig())

	storage.config.logLevel = "debug"
	storage.config.storageDir = "./"

	err := storage.SessionStart(context.Background())
	assert.Nil(t, err)
}
