package main

import (
	"fmt"
	"os"
	"time"

	// external
	pp "github.com/k0kubun/pp"

	// internal
	storage "github.com/sniperkit/colly-storage/pkg"
	dal_pivot "github.com/sniperkit/colly-storage/plugin/dal/pivot"
)

var (
	storageDebug        bool          = false
	storageBackend      string        = "mysql" // use a flag to switch between backends
	storagePingDuration time.Duration = 5 * time.Second
)

var (
	store         storage.Storage
	storageConfig *dal_pivot.Config
)

func main() {

	fmt.Println("Running storage `colly-dal-pivot` example...")

	var err error

	switch storageBackend {
	case "sqlite":
		storageConfig = &dal_pivot.Config{
			Scheme:  "sqlite",                              // required
			Host:    "",                                    // required
			Dataset: "./shared/storage/dal/pivot/colly.db", // required
			Options: map[string]interface{}{},              // optional
		}
	case "mysql":
		storageConfig = &dal_pivot.Config{
			DSN: `mysql://test:test@localhost:3306/colly?parseTime=True&loc=Local`,
			// mysql://test:test@tcp(localhost:3306)/colly-data/qeue?parseTime=True&loc=Local`,
		}

	case "postgres":
		storageConfig = &dal_pivot.Config{
			Scheme:  "postgres",                                   // required
			Host:    "test:test@localhost:5432",                   // required
			Dataset: "colly",                                      // required
			Options: map[string]interface{}{"sslmode": "disable"}, // optional
		}

	case "mongodb":
		storageConfig = &dal_pivot.Config{
			DSN: `mongodb://localhost:27017/colly`,
		}

	case "tiedot":
		storageConfig = &dal_pivot.Config{
			DSN: `tiedot://colly`,
		}

	}

	store, err = dal_pivot.NewDataAbstractionLayer(storageConfig)
	if err != nil {
		fmt.Println("error while creating a new data abstraction layer instance... error=", err)
		os.Exit(1)
	}

	if storageDebug {
		pp.Println("Storage=", store)
	} else {
		store.Action("ping", storagePingDuration)
		store.Action("list_collections", nil)
	}

}

/*
func createCollection(name string) error {
	return store.CreateCollection(
		store.NewCollection(name),
	)
}

func checkCollection(name string) error {
	if coll, err := store.GetCollection(name); err != nil {
		return err
	}
	return nil
}

func loadData(collection string) error {

	// Create collection with schema
	// --------------------------------------------------------------------------------------------
	err := store.CreateCollection(
		store.NewCollection(collection).
			AddFields(store.Field{
				Name: `name`,
				Type: store.StringType,
			}, store.Field{
				Name:         `created_at`,
				Type:         store.TimeType,
				DefaultValue: time.Now,
			}))

	var record *store.Record

	// Insert and Retrieve
	// --------------------------------------------------------------------------------------------
	recordset := store.NewRecordSet(
		store.NewRecord(testCrudIdSet[0]).Set(`name`, `First`),
		store.NewRecord(testCrudIdSet[1]).Set(`name`, `Second`),
		store.NewRecord(testCrudIdSet[2]).Set(`name`, `Third`))
}

func search(collection string) error {
	collection := store.NewCollection(`TestSearchQuery`).
		AddFields(store.Field{
			Name: `name`,
			Type: store.StringType,
		})

	// set the global page size at the package level for this test
	backends.IndexerPageSize = 5
	backends.IndexerPageSize = 100
	if search := store.WithSearch(collection); search != nil {
		err := store.CreateCollection(collection)

		// f, err := filter.Parse(`all`)
		// f.Offset = 20
		// f.Limit = 9
		// keyValues, err := search.ListValues(collection, []string{`name`}, filter.All())
		// v, ok := keyValues[`name`]
	}
	// if agg := backend.WithAggregator(collection); agg != nil {
	// }
}
*/
