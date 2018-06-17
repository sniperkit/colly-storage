package badgerstorage

import (
	"errors"
)

func (s *Store) Init() error {
	return errors.New("Init method is not implemented yet...")
}

func (s *Store) Debug(action string) error {
	// s.listAll()
	// s.keys()
	// s.compressor()
	// s.db.PurgeOlderVersions()
	// s.purgeOlderVersions()
	// s.updates()
	// s.seekPrefix()
	return errors.New("Debug is not implemented yet")
}

func (s *Store) Clear() error { return errors.New("Debug is not implemented yet") }

// Ping connects to the database. Returns nil if successful.
func (s *Store) Ping() error {
	return errors.New("Ping() method is not implemented yet")
}
