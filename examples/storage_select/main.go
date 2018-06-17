package main

import (
	"fmt"
	"os"

	// external
	pp "github.com/sniperkit/pp"

	// internal
	storage "github.com/sniperkit/colly-storage/pkg"
	bck_badger "github.com/sniperkit/colly-storage/plugin/backend/badger"
	bck_boltdb "github.com/sniperkit/colly-storage/plugin/backend/boltdb"
	bck_bboltdb "github.com/sniperkit/colly-storage/plugin/backend/boltdb_bbolt"
)

var (
	backendName    string = "badger" // use a flag to switch between backends
	backendDefault string = "in_memory"
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
		store, err = bck_badger.New(&bck_badger.Config{})

	case "boltdb":
		store, err = bck_boltdb.New(&bck_boltdb.Config{})

	case "bbolt":
		store, err = bck_bboltdb.New(&bck_bboltdb.Config{})

	case "in_memory":
		fallthrough

	default:
		store, err = storage.New(&storage.InMemoryConfig{MaxRow: 100000})

	}

	if err != nil {
		fmt.Println("error while creating a new storage instance...")
		os.Exit(1)
	}

	pp.Println("Storage", store)

}
