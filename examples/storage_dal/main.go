package main

import (
	"fmt"
	"os"

	// external
	pp "github.com/k0kubun/pp"

	// internal
	storage "github.com/sniperkit/colly-storage/pkg"
	bck_dal "github.com/sniperkit/colly-storage/plugin/dal"
)

var (
	storageBackendDSN string = bck_dal.DefaultBackendDSN
	//
	// Backend DSN
	// ref. https://github.com/ghetzel/pivot/blob/master/backends/backends.go#L41
	//
	// sqlite3: 		`sqlite:///./tmp/db_test/test.db`
	// mysql/mariadb: 	`mysql://test:test@db/test`
	// dynamoDB: 		`dynamodb://test:test@db/test`
	// postgres: 		`postgres://test:test@db/test`
	// file:			`file://%s/`
	// filesystem: 		`fs://%s/`
	// filesystem+yaml: `fs+yaml://%s/`
	// filesystem+json: `fs+json://%s/`
	// tiedot: 			`tiedot://%s/`
	// mongodb: 		`mongodb://localhost/test`
	// elastic: 		`not ready yet`
	storagePrefixPath string = bck_dal.DefaultStoragePrefixPath
)

func main() {

	fmt.Println("Running storage `colly-dal` example...")

	var err error
	var store storage.Storage

	fmt.Println("Starting `colly-dal` storage, dsn=", storageBackendDSN, ", prefixPath", storagePrefixPath)
	conf := &bck_dal.Config{
		PrefixPath: &storageBackendDSN,
		DSN:        storagePrefixPath,
		Verbose:    true,
		Debug:      false,
		Sanitize:   false,
	}
	store, err = bck_dal.NewDataAbstractionLayer(conf)
	if err != nil {
		fmt.Println("error while creating a new data abstraction layer instance... error=", err)
		os.Exit(1)
	}

	pp.Println("Storage=", store)

}
