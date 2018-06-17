package storage

import (
	"errors"
	"sync"
)

// InMemoryStorage is the default storage backend of storage interface.
type InMemoryStorage struct {
	lock *sync.RWMutex
	// visitedURLs map[uint64]bool
	// jar         *cookiejar.Jar
}

type InMemoryConfig struct {
	MaxRow int
}

// Init initializes InMemoryStorage
func (s *InMemoryStorage) Init() error {
	return errors.New("Init method is not implemented yet...")
}

func (s *InMemoryStorage) Get(key string) (resp []byte, ok bool) {
	return []byte{}, false
}

func (s *InMemoryStorage) Set(key string, resp []byte) error {
	return errors.New("Set() method is not implemented yet")
}

func (s *InMemoryStorage) Delete(key string) error {
	return errors.New("Delete() method is not implemented yet")
}

// Debug
func (s *InMemoryStorage) Debug(action string) error {
	return errors.New("Debug() method is not implemented yet")
}

// Ping check if the storage is available...
func (s *InMemoryStorage) Ping() error {
	return errors.New("Ping() method is not implemented yet")
}

// Clear truncate all key/values stored...
func (s *InMemoryStorage) Clear() error {
	return errors.New("Debug is not implemented yet")
}

// Close deletes the storage
func (s *InMemoryStorage) Close() error {
	return nil
}

/*
// Visited implements Storage.Visited()
func (s *InMemoryStorage) Visited(requestID uint64) error {
	s.lock.Lock()
	s.visitedURLs[requestID] = true
	s.lock.Unlock()
	return nil
}

// IsVisited implements Storage.IsVisited()
func (s *InMemoryStorage) IsVisited(requestID uint64) (bool, error) {
	s.lock.RLock()
	visited := s.visitedURLs[requestID]
	s.lock.RUnlock()
	return visited, nil
}

// Cookies implements Storage.Cookies()
func (s *InMemoryStorage) Cookies(u *url.URL) string {
	return StringifyCookies(s.jar.Cookies(u))
}

// SetCookies implements Storage.SetCookies()
func (s *InMemoryStorage) SetCookies(u *url.URL, cookies string) {
	s.jar.SetCookies(u, UnstringifyCookies(cookies))
}
*/
