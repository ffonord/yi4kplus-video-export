package localdisk

import (
	"context"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain/file"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestStorage_SessionStart(t *testing.T) {
	storage := New(NewConfig(), logger.New(logger.EnvTest))

	storage.config.storageDir = "./"

	err := storage.SessionStart(context.Background())
	assert.Nil(t, err)
}

func TestStorage_GetWriter(t *testing.T) {
	storage := New(NewConfig(), logger.New(logger.EnvTest))

	storage.config.storageDir = "."

	filePath := "./test_write.file"
	fileContent := "some content"

	f := file.New(
		"test_write.file",
		"./",
		time.Now(),
		uint64(len(fileContent)),
	)

	wc, err := storage.GetWriter(f)
	assert.Nil(t, err)

	w, err := wc.Write([]byte(fileContent))

	assert.Nil(t, err)
	assert.Equal(t, len(fileContent), w)

	err = os.Remove(filePath)
	assert.Nil(t, err)
}

func TestStorage_Delete(t *testing.T) {
	storage := New(NewConfig(), logger.New(logger.EnvTest))

	storage.config.storageDir = "."

	filePath := "./test_delete.file"
	fileContent := "some content"

	_ = os.WriteFile(filePath, []byte(fileContent), 0644)

	f := file.New(
		"test_delete.file",
		"./",
		time.Now(),
		uint64(len(fileContent)),
	)

	err := storage.Delete(f)
	assert.Nil(t, err)
}
