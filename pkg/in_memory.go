package storage

import (
	"errors"
	"sync"
)

// Store is the default storage backend of storage interface.
type Store struct {
	lock *sync.RWMutex
	// db          map[uint64]bool
	// visitedURLs map[uint64]bool
	// jar         *cookiejar.Jar
}

type Config struct {
	MaxRow     int
	Sanitize   bool
	ReadOnly   bool
	StrictMode bool
	Compress   bool
	Debug      bool
	Stats      bool
}

func NewInMemoryStorage(config *Config) (*Store, error) {
	s := &Store{
		lock: &sync.RWMutex{},
	}
	return s, nil
}

// func Mount() *Store {
//	return &Store{}
//}

// Init initializes Store
func (s *Store) Init() error {
	return errors.New("Init method is not implemented yet...")
}

func (s *Store) Get(key string) (resp []byte, ok bool) {
	return []byte{}, false
}

func (s *Store) Set(key string, resp []byte) error {
	return errors.New("Set() method is not implemented yet")
}

func (s *Store) Delete(key string) error {
	return errors.New("Delete() method is not implemented yet")
}

// Debug
func (s *Store) Debug(action string) error {
	return errors.New("Debug() method is not implemented yet")
}

func (c *Store) Action(name string, args ...interface{}) (map[string]*interface{}, error) {
	return nil, errors.New("Action() method is not implemented yet")
}

// Ping check if the storage is available...
func (s *Store) Ping() error {
	return errors.New("Ping() method is not implemented yet")
}

// Clear truncate all key/values stored...
func (s *Store) Clear() error {
	return errors.New("Debug is not implemented yet")
}

// Close deletes the storage
func (s *Store) Close() error {
	return nil
}

/*
// Visited implements Storage.Visited()
func (s *Store) Visited(requestID uint64) error {
	s.lock.Lock()
	s.visitedURLs[requestID] = true
	s.lock.Unlock()
	return nil
}

// IsVisited implements Storage.IsVisited()
func (s *Store) IsVisited(requestID uint64) (bool, error) {
	s.lock.RLock()
	visited := s.visitedURLs[requestID]
	s.lock.RUnlock()
	return visited, nil
}

// Cookies implements Storage.Cookies()
func (s *Store) Cookies(u *url.URL) string {
	return StringifyCookies(s.jar.Cookies(u))
}

// SetCookies implements Storage.SetCookies()
func (s *Store) SetCookies(u *url.URL, cookies string) {
	s.jar.SetCookies(u, UnstringifyCookies(cookies))
}
*/
