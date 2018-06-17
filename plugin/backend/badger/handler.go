package badgerstorage

import (
	"errors"
	"path/filepath"
	"sync"
	"time"

	// external
	"github.com/dgraph-io/badger"
)

var defaultStorePath string = filepath.Join(defaultStoreDir, defaultStoreValueDir)

// Store stores and retrieves data using Badger KV.
type Store struct {
	mu syns.RWMutex
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


// ListAll lists all the pairs KV of a given type
func (s *Store) compressor() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var key string
	in := []byte(`HTTP/1.1 200 OK\r\nTransfer-Encoding: chunked\r\nAccess-Control-Allow-Origin: *\r\nAccess-Control-Expose-Headers: ETag, Link, Retry-After, X-GitHub-OTP, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, X-OAuth-Scopes, X-Accepted-OAuth-Scopes, X-Poll-Interval\r\nStore-Control: private, max-age=60, s-maxage=60\r\nContent-Security-Policy: default-src 'none'\r\nContent-Type: application/json; charset=utf-8\r\nDate: Sun, 17 Dec 2017 09:01:26 GMT\r\nEtag: W/\"97327ebe7afdfe040e366a8f35cf8c72\"\r\nLast-Modified: Fri, 26 May 2017 20:57:47 GMT\r\nServer: GitHub.com\r\nStatus: 200 OK\r\nStrict-Transport-Security: max-age=31536000; includeSubdomains; preload\r\nVary: Accept, Authorization, Cookie, X-GitHub-OTP\r\nX-Accepted-Oauth-Scopes: \r\nX-Content-Type-Options: nosniff\r\nX-Frame-Options: deny\r\nX-Github-Media-Type: github.v3; format=json\r\nX-Github-Request-Id: EA31:2902B:1281CE8:23CD027:5A363265\r\nX-Oauth-Scopes: repo, user\r\nX-Ratelimit-Limit: 5000\r\nX-Ratelimit-Remaining: 2441\r\nX-Ratelimit-Reset: 1513504311\r\nX-Runtime-Rack: 0.049262\r\nX-Varied-Accept: application/vnd.github.v3+json\r\nX-Varied-Authorization: Bearer 63814c0ef8a9a7a273e828d1cc4d410b4f449a9f\r\nX-Xss-Protection: 1; mode=block\r\n\r\n347\r\n{\"name\":\"README.md\",\"path\":\"README.md\",\"sha\":\"723b022f9b66c690b95239ef7de83f9dc9d24290\",\"size\":60,\"url\":\"https://api.github.com/repos/AaronTL/TapNews/contents/README.md?ref=master\",\"html_url\":\"https://github.com/AaronTL/TapNews/blob/master/README.md\",\"git_url\":\"https://api.github.com/repos/AaronTL/TapNews/git/blobs/723b022f9b66c690b95239ef7de83f9dc9d24290\",\"download_url\":\"https://raw.githubusercontent.com/AaronTL/TapNews/master/README.md\",\"type\":\"file\",\"content\":\"IyBUYXBOZXdzClJlYWwgVGltZSBOZXdzIFNjcmFwaW5nIGFuZCBSZWNvbW1l\\nbmRhdGlvbiBTeXN0ZW0K\\n\",\"encoding\":\"base64\",\"_links\":{\"self\":\"https://api.github.com/repos/AaronTL/TapNews/contents/README.md?ref=master\",\"git\":\"https://api.github.com/repos/AaronTL/TapNews/git/blobs/723b022f9b66c690b95239ef7de83f9dc9d24290\",\"html\":\"https://github.com/AaronTL/TapNews/blob/master/README.md\"}}\r\n0\r\n\r\n`)

	key = "test_raw"
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), in)
	})
	if err != nil {
		return err
	}
	Match("application/json", in)
	key = "test_gzipped"
	return s.db.Update(func(txn *badger.Txn) error {
		in, err = gzipData(in)
		if err != nil {
			return err
		}
		return txn.Set([]byte(key), in)
	})
}

// ListAll lists all the pairs KV of a given type
func (s *Store) listAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		i := 0
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			kStr := string(k)
			v, err := item.Value()
			if err != nil {
				return err
			}
			vStr := string(v)

			ling := DetectLang(vStr)
			lang, safe := DetectType(string(k), vStr)
			i++
		}
		return nil
	})
	return err
}

func (s *Store) keys() (keys []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.View(func(txn *badger.Txn) error {
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
}

func (s *Store) purgeOlderVersions() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.PurgeOlderVersions()
}

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

// Add a hashed key to the store if it doesn't already exist
func (s *Store) touchHashed(in string) (err error) {
	return s.db.Touch([]byte(roencoding.XXHash(in)))
}

// Does hash of key exist in the store?
func (s *Store) existsHashed(in string) (exists bool, err error) {
	return s.db.Exists([]byte(roencoding.XXHash(in)))
}

/*
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
