package bboltstorage

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	// external
	bbolt "github.com/coreos/bbolt"

	// internal
	helper "github.com/sniperkit/colly-storage/pkg/helper"
)

var (
	DefaultStorageDir  string = filepath.Join(StoragePrefixPath, StorageBucketName)
	DefaultStorageFile string = fmt.Sprintf("%s/%s%s", StoragePrefixPath, StorageBucketName, StorageFileExtension)
)

// Store is an implementation of httpstore.Store that uses a bolt database.
type Store struct {
	// sync.Mutex
	sync.RWMutex
	db          *bbolt.DB
	storagePath string
	bucketName  string
	fp          string
	debug       bool
	stats       bool
	compress    bool
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

// New returns a new Store that uses a bolt database at the given path.
func New(config *Config) (*Store, error) {

	if config == nil {
		config.StoragePath = DefaultStorageFile
		config.BucketName = StorageBucketName
	}

	if config.StoragePath == "" {
		config.StoragePath = DefaultStorageFile
	}

	if config.BucketName == "" {
		config.BucketName = StorageBucketName
	}

	var err error
	store := &Store{}
	store.storagePath = config.StoragePath
	store.bucketName = config.BucketName
	store.compress = config.Compress
	store.debug = config.Debug
	store.stats = config.Stats

	if err := helper.EnsurePathExists(config.StoragePath); err != nil {
		return nil, err
	} else {
		store.fp = fmt.Sprintf("%s/%s%s", config.StoragePath, StorageBucketName, StorageFileExtension)
	}

	store.db, err = bbolt.Open(config.StoragePath, 0755, nil)
	if err != nil {
		return nil, err
	}

	init := func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(config.BucketName))
		return err
	}

	if err := store.db.Update(init); err != nil {
		return nil, err
		if err := store.db.Close(); err != nil {
			return nil, err
		}
	}
	return store, nil
}

// Mount returns a new Store using the provided (and opened) bolt database.
func Mount(db *bbolt.DB) *Store {
	return &Store{db: db}
}

// Close closes the underlying boltdb database.
func (c *Store) Close() error {
	return c.db.Close()
}

// Get retrieves the response corresponding to the given key if present.
func (c *Store) Get(key string) (resp []byte, ok bool) {
	c.RLock()
	get := func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(c.bucketName))
		if bkt == nil {
			return errors.New("bucket is nil")
		}
		resp = bkt.Get([]byte(key))
		return nil
	}
	if err := c.db.View(get); err != nil {
		return resp, false
	}
	c.RUnlock()
	if c.compress {
		var err error
		resp, err = ungzipData(resp)
		if err != nil {
			return resp, false
		}
	}
	return resp, resp != nil
}

// Set stores a response to the store at the given key.
func (c *Store) Set(key string, resp []byte) error {
	c.Lock()
	set := func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(c.bucketName))
		if bkt == nil {
			return errors.New("bucket is nil")
		}
		if c.compress {
			var err error
			resp, err = gzipData(resp)
			if err != nil {
				return errors.New("error while compressing content...")
			}
		}
		return bkt.Put([]byte(key), resp)
	}
	c.Unlock()
	if err := c.db.Update(set); err != nil {
		return err
	}
	return nil
}

// Delete removes the response with the given key from the store.
func (c *Store) Delete(key string) error {
	c.Lock()
	del := func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(c.bucketName))
		if bkt == nil {
			return errors.New(fmt.Sprintf("bboltstore.Delete(): could not reach the bucket: %s", c.bucketName))
		}
		return bkt.Delete([]byte(key))
	}
	if err := c.db.Update(del); err != nil {
		return err
	}
	c.Unlock()
	return nil
}

func ungzipData(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err = ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func gzipData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func fmtToJsonArr(s []byte) []byte {
	s = bytes.Replace(s, []byte("{"), []byte("[{"), 1)
	s = bytes.Replace(s, []byte("}"), []byte("},"), -1)
	s = bytes.TrimSuffix(s, []byte(","))
	s = append(s, []byte("]")...)
	return s
}
