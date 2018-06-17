package badgerstorage

import (
	"errors"
	"path/filepath"
	"sync"
	"time"

	// external
	"github.com/dgraph-io/badger"
	"github.com/rohanthewiz/roencoding"
)

var defaultStorePrefixPath string = filepath.Join(storagePrefixPath, storageBucketName)

// Store stores and retrieves data using Badger KV.
type Store struct {
	mu          sync.RWMutex
	db          *badger.DB
	storagePath string
	bucketName  string
	compress    bool
	debug       bool
}

type Check struct {
	Enabled   bool
	Key       string
	Requests  int
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
	Priority  bool
	Provider  string
}

func Mount(client *badger.DB) *Store {
	return &Store{db: client}
}

func New(config *Config) (*Store, error) {
	badgerConfig := badger.DefaultOptions
	if config == nil {
		badgerConfig.Dir = storagePrefixPath
		badgerConfig.ValueDir = storageBucketName
		badgerConfig.SyncWrites = true
	} else {
		badgerConfig.Dir = config.StoragePath
		badgerConfig.ValueDir = filepath.Join(config.StoragePath, config.ValueDir)
		badgerConfig.SyncWrites = config.SyncWrites
	}

	client, err := badger.Open(badgerConfig)
	if err != nil {
		return nil, err
	}

	return &Store{
		db:       client,
		debug:    config.Debug,
		compress: config.Compress,
	}, nil
}

func (s *Store) Get(key string) (resp []byte, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		resp, err = item.Value()
		if err != nil {
			return err
		}
		if s.compress {
			var err error
			resp, err = Decompress(resp)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return resp, err == nil
}

// Set stores a response to the store at the given key.
func (s *Store) Set(key string, resp []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Update(func(txn *badger.Txn) error {
		if s.compress {
			var err error
			resp, err = Compress(resp)
			if err != nil {
				return errors.New("error while compressing content...")
			}
		}
		return txn.Set([]byte(key), resp)
	})
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
}

func (s *Store) Action(name string, args ...interface{}) (map[string]*interface{}, error) {
	if name == "getKeys" {
		keys := s.keys()
		resp := make(map[string]*interface{})
		for _, v := range keys {
			resp[v] = nil
		}
		return resp, nil
	}

	return nil, errors.New("Action not implemented yet")
}

// Close closes the underlying boltdb database.
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) keys() (keys []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			keys = append(keys, string(k))
		}
		return nil
	})
	if err != nil {
		return []string{}
	}
	return
}

// Add a hashed key to the store if it doesn't already exist
func (s *Store) touchHashed(key string) error {
	return s.Set(roencoding.XXHash(key), []byte(roencoding.XXHash(key)))
}

// Does hash of key exist in the store?
func (s *Store) existsHashed(key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(roencoding.XXHash(key)))
		if err != nil {
			return err
		}
		return nil
	})
	return err == nil, err
}

/*

func (s *Store) seekPrefix(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		prefix := []byte(value)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) firstIndex() (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx := s.db.NewTransaction(false)
	defer tx.Discard()
	iter := tx.NewIterator(iterAscOpt)
	iter.Rewind()
	item := iter.Item()
	if item == nil {
		return 0, nil
	}

	return bytesToUint64(item.Key()), nil
}

func (s *Store) lastIndex() (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx := s.db.NewTransaction(false)
	defer tx.Discard()
	iter := tx.NewIterator(iterDescOpt)
	iter.Rewind()
	item := iter.Item()
	if item == nil {
		return 0, nil
	}
	return bytesToUint64(item.Key()), nil
}

func (s *Store) deleteRange(min, max uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx := s.db.NewTransaction(true)
	defer tx.Discard()
	minKey := uint64ToBytes(min)
	iter := tx.NewIterator(iterAscOpt)
	for iter.Seek(minKey); iter.Valid(); iter.Next() {
		item := iter.Item()
		if item == nil {
			break
		}
		curKey := safeKey(item)
		if bytesToUint64(curKey) > max {
			break
		}
		if err := tx.Delete(curKey); err != nil {
			return err
		}
	}
	if err := tx.Commit(nil); err != nil {
		return err
	}
	return nil
}
*/

/*
func (s *Store) updates(updates map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	txn := db.NewTransaction(true)
	for k,v := range updates {
	  if err := txn.Set(byte[](k),byte[](v)); err == ErrTxnTooBig {
	    _ = txn.Commit()
	    txn = db.NewTransaction(..)
	    _ = txn.Set(k,v)
	  }
	}
	_ = txn.Commit()

}
*/
