package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFileName = "test"

func TestStore(t *testing.T) {
	assert.NoError(t, clearCache(testFileName))

	store := New(testFileName)
	assert.NotNil(t, store)

	store.Set("key1", "foo")

	assert.Equal(t, "foo", store.Get("key1"))

	assert.NoError(t, store.Flush())

	store = New(testFileName)
	assert.NotNil(t, store)

	assert.Equal(t, "foo", store.Get("key1"))
}

func TestCorruptedFile(t *testing.T) {
	assert.NoError(t, clearCache(testFileName))

	store := New(testFileName)
	assert.NotNil(t, store)

	store.Set("key1", "foo")
	store.Flush()

	// Corrupt the file
	filepath, err := getCacheFilePath(testFileName)
	assert.NoError(t, err)
	f, err := os.OpenFile(filepath, os.O_WRONLY, 0)
	assert.NoError(t, err)
	defer f.Close()
	_, err = f.Write([]byte("corrupted"))
	assert.NoError(t, err)
	f.Close()

	store = New(testFileName)
	assert.NotNil(t, store)

	assert.Nil(t, store.Get("key1"))
}
