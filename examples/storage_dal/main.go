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
	backendEngine     string = "sqlite3" // use a flag to switch between backends
	backendDefault    string = "sqlite3"
	storagePrefixPath string = "./shared/"
)

func main() {

	fmt.Println("Running storage `colly-dal` example...")

	var err error
	var store storage.Storage

	if backendEngine == "" {
		backendEngine = backendDefault
	}

	fmt.Println("Starting storage, backendEngine=", backendEngine)
	conf := &bck_dal.Config{
		Engine: backendEngine,
		DSN:    "./shared/storage/sqlite3/colly-dal.sqlite",
	}
	store, err = bck_dal.NewDataAbstractionLayer(conf)
	if err != nil {
		fmt.Println("error while creating a new data abstraction layer instance... error=", err)
		os.Exit(1)
	}

	pp.Println("Storage=", store)

}
