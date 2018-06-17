package bboltstorage

import (
	"errors"
)

func (s *Store) Init() error {
	return errors.New("Init method is not implemented yet...")
}

func (c *Store) Action(name string, args ...interface{}) (map[string]*interface{}, error) {
	return nil, errors.New("Action() method is not implemented yet")
}

// Ping connects to the database. Returns nil if successful.
func (s *Store) Ping() error {
	return errors.New("Ping() method is not implemented yet")
}

func (s *Store) Debug(action string) error { return errors.New("Debug() method is not implemented yet") }

func (s *Store) Clear() error { return errors.New("Clear() method is not implemented yet") }
