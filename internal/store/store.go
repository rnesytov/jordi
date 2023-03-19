package store

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/adrg/xdg"
	"github.com/pkg/errors"
)

const appName = "jordi"

type (
	Value any

	Store struct {
		lock     sync.RWMutex
		filepath string
		data     map[string]Value
	}
)

func New(name string) *Store {
	cacheFilePath, _ := getCacheFilePath(name)

	store := &Store{
		filepath: cacheFilePath,
		data:     make(map[string]Value),
	}
	if err := store.load(); err != nil {
		fmt.Println(err)
	}

	return store
}

func (s *Store) load() error {
	file, err := os.Open(s.filepath)

	if err == nil {
		defer file.Close()

		if err := json.NewDecoder(file).Decode(&s.data); err != nil {
			return errors.Wrapf(err, "failed to decode cache file '%s'", s.filepath)
		}
	}
	return nil
}

func (s *Store) Get(key string) Value {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.data[key]
}

func (s *Store) Set(key string, value Value) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
}

func (s *Store) Flush() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	file, err := os.Create(s.filepath)

	if err == nil {
		defer file.Close()

		if err := json.NewEncoder(file).Encode(&s.data); err != nil {
			return errors.Wrapf(err, "failed to encode cache file '%s'", s.filepath)
		}
	}
	return nil
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func getCacheFilePath(name string) (string, error) {
	return xdg.CacheFile(fmt.Sprintf("%s/%s.json", appName, md5Hash(name)))
}

func clearCache(name string) error {
	cacheFilePath, err := getCacheFilePath(name)
	if err != nil {
		return err
	}
	os.Remove(cacheFilePath)
	return nil
}
