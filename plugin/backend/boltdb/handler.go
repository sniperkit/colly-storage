package boltdbstorage

import (
	"errors"
	// "path/filepath"
	"fmt"
	"sync"

	// external
	"github.com/boltdb/bolt"

	// internal
	helper "github.com/sniperkit/colly-storage/pkg/helper"
)

var (
	DefaultStorageFile string = fmt.Sprintf("%s/%s/%s%s", StoragePrefixPath, StorageBucketName, StorageBucketName, StorageFileExtension)
)

type Store struct {
	sync.RWMutex
	db          *bolt.DB
	storagePath string
	bucketName  string
	fp          string
	debug       bool
}

func New(config *Config) (*Store, error) {
	if config == nil {
		config.StoragePath = DefaultStorageFile
	}
	if config.StoragePath == "" {
		return nil, errors.New("boltdbstore.New(): Storage path is not defined.")
	}
	if config.BucketName == "" {
		config.BucketName = StorageBucketName
	}

	store := &Store{}
	store.debug = config.Debug

	if err := helper.EnsurePathExists(config.StoragePath); err != nil {
		return nil, err
	} else {
		store.fp = fmt.Sprintf("%s/%s%s", config.StoragePath, StorageBucketName, StorageFileExtension)
	}

	var err error
	store.db, err = bolt.Open(config.StoragePath, 0600, nil)
	if err != nil {
		return nil, err
	}

	init := func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(config.BucketName))
		return err
	}

	if err := store.db.Update(init); err != nil {
		if err := store.db.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}
	return store, nil
}

// Mount returns a new Cache using the provided (and opened) bolt database.
func Mount(db *bolt.DB) *Store {
	return &Store{db: db}
}

// Close closes the underlying boltdb database.
func (s *Store) Close() error {
	return s.db.Close()
}

// Get retrieves the response corresponding to the given key if present.
func (s *Store) Get(key string) (resp []byte, ok bool) {
	s.RLock()

	get := func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(s.bucketName))
		if bkt == nil {
			return errors.New("bucket is nil")
		}
		resp = bkt.Get([]byte(key))
		return nil
	}
	if err := s.db.View(get); err != nil {
		return resp, false
	}
	s.RUnlock()
	return resp, resp != nil
}

// Set stores a response to the store at the given key.
func (s *Store) Set(key string, resp []byte) error {
	s.Lock()

	set := func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(s.bucketName))
		if bkt == nil {
			return errors.New("bucket is nil")
		}
		return bkt.Put([]byte(key), resp)
	}
	s.db.Update(set)
	if err := s.db.Update(set); err != nil {
		return err
	}
	s.Unlock()
	return nil
}

// Delete removes the response with the given key from the store.
func (s *Store) Delete(key string) error {
	s.Lock()

	del := func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(s.bucketName))
		if bkt == nil {
			return errors.New("bucket is nil")
		}
		return bkt.Delete([]byte(key))
	}
	if err := s.db.Update(del); err != nil {
		return err
	}
	s.Unlock()
	return nil
}

// Ping connects to the database. Returns nil if successful.
func (s *Store) Ping() error {
	return s.db.View(func(tx *bolt.Tx) error { return nil })
}
