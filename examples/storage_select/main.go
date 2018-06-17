package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	// external
	pp "github.com/k0kubun/pp"

	// internal
	storage "github.com/sniperkit/colly-storage/pkg"
	bck_badger "github.com/sniperkit/colly-storage/plugin/backend/badger"
	bck_badger_lru "github.com/sniperkit/colly-storage/plugin/backend/badger_lru"
	bck_boltdb "github.com/sniperkit/colly-storage/plugin/backend/boltdb"
	bck_bboltdb "github.com/sniperkit/colly-storage/plugin/backend/boltdb_bbolt"
)

var (
	backendName       string = "bbolt" // use a flag to switch between backends
	backendDefault    string = "in_memory"
	storagePrefixPath string = "./shared/"
)

func main() {

	fmt.Println("Running storage backend select example...")

	var err error
	var store storage.Storage

	if backendName == "" {
		backendName = backendDefault
	}

	fmt.Println("Starting storage, backend=", backendName)

	switch backendName {
	case "badger":
		// conf
		conf := &bck_badger.Config{
			ValueDir:    "colly-storage.snappy",
			StoragePath: filepath.Join(storagePrefixPath, "storage", "badger"),
			SyncWrites:  false,
			Debug:       false,
			Compress:    true,
			TTL:         time.Duration(120 * 24 * time.Hour),
		}
		// init
		store, err = bck_badger.New(conf)

	case "badger-lru":
		// conf
		conf := &bck_badger_lru.Config{
			ValueDir:    "colly-storage",
			StoragePath: filepath.Join(storagePrefixPath, "storage", "badger-lru"),
			SyncWrites:  false,
			Debug:       false,
			Compress:    true,
			TTL:         time.Duration(120 * 24 * time.Hour),
		}
		// init
		store, err = bck_badger_lru.New(conf)

	case "boltdb":
		// conf
		bucketName := "colly-storage"
		bucketPrefixPath := filepath.Join(storagePrefixPath, "storage", "boltdb")
		bucketStorageFilePath := fmt.Sprintf("%s/%s/%s%s", bucketPrefixPath, bucketName, bucketName, bck_boltdb.StorageFileExtension)
		conf := &bck_boltdb.Config{
			BucketName:  "colly-storage.boltdb",
			StoragePath: bucketStorageFilePath,
			Debug:       false,
		}
		// init
		store, err = bck_boltdb.New(conf)

	case "bbolt":
		// conf
		bucketName := "colly-storage"
		bucketPrefixPath := filepath.Join(storagePrefixPath, "storage", "bbolt")
		bucketStorageFilePath := fmt.Sprintf("%s/%s/%s%s", bucketPrefixPath, bucketName, bucketName, bck_bboltdb.StorageFileExtension)
		conf := &bck_bboltdb.Config{
			BucketName:  "colly-storage.bbolt",
			StoragePath: bucketStorageFilePath,
			Debug:       false,
			Stats:       false,
		}
		// init
		store, err = bck_bboltdb.New(conf)

	case "diskv":
		// cacheStoragePrefixPath := filepath.Join(prefixPath, "cacher.diskv")
		// fsutil.EnsureDir(cacheStoragePrefixPath)
		// conf := &bck_diskv.Config{}
		// store, err = diskcache.New(conf)

	case "in_memory":
		fallthrough

	default:
		conf := &storage.Config{
			MaxRow: 100000,
			Debug:  false,
			Stats:  false,
		}
		store, err = storage.NewInMemoryStorage(conf)

	}

	if err != nil {
		fmt.Println("error while creating a new storage instance... error=", err)
		os.Exit(1)
	}

	pp.Println("Storage", store)

}
