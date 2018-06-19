package dal_pivot

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	// external
	"github.com/ghetzel/pivot"
	"github.com/ghetzel/pivot/backends"
	"github.com/ghetzel/pivot/dal"
	"github.com/ghetzel/pivot/mapper"

	// pp "github.com/k0kubun/pp"

	// internal
	helper "github.com/sniperkit/colly-storage/pkg/helper"
)

var (
	DefaultBackendDSN string = fmt.Sprintf(`sqlite:///%s/%s%s`, DefaultStoragePrefixPath, DefaultStorageDatabaseName, DefaultStorageFileExtension)
)

type Store struct {
	lock    *sync.RWMutex
	backend backends.Backend
	mapper  mapper.Mapper
	schema  []*dal.Collection
	conf    *Config
}

// alias
func MakeConnectionString(scheme string, host string, dataset string, options map[string]interface{}) (dal.ConnectionString, error) {
	return dal.MakeConnectionString(scheme, host, dataset, options)
}

func checkMissingParameter(params map[string]string) (missing []string) {
	for k, v := range params {
		if v == "" {
			missing = append(missing, k)
		}
	}
	return
}

func NewDataAbstractionLayer(config *Config) (*Store, error) {
	s := &Store{
		lock: &sync.RWMutex{},
	}

	if config.DSN == "" {
		if config.Scheme == "sqlite" || config.Scheme == "sqlite3" || config.Scheme == "boltdb" || config.Scheme == "badger" || config.Scheme == "tiedot" {
			if err := helper.EnsurePathExists(config.Dataset); err != nil {
				return nil, err
			}
		} else {
			p := map[string]string{
				"scheme":  config.Scheme,
				"host":    config.Host,
				"dataset": config.Dataset,
			}
			if missingKeys := checkMissingParameter(p); len(missingKeys) > 0 {
				return nil, errors.New(fmt.Sprintf("Empty values for required parameters: %s", strings.Join(missingKeys, ",")))
			}
		}
		if conn, err := dal.MakeConnectionString(config.Scheme, config.Host, config.Dataset, config.Options); err != nil {
			return nil, err
		} else {
			config.DSN = conn.String()
		}
	}

	// setup a new backend instance based on the supplied connection string
	backend, err := pivot.NewDatabase(config.DSN)
	if err != nil {
		return nil, err
	}
	s.backend = backend

	// initialize the backend (connect to/open it)
	if err := s.backend.Initialize(); err != nil {
		return nil, err
	}

	// copy config
	s.conf = config

	return s, nil
}

func (c *Store) Action(name string, args ...interface{}) (map[string]*interface{}, error) {
	switch name {
	case "ping":
		c.backend.Ping(args[0].(time.Duration))

	case "getCollections":
		c.backend.ListCollections()

	}

	return nil, errors.New("Action() method is not implemented yet")
}

func (s *Store) Get(key string) (resp []byte, ok bool) {
	return []byte{}, false
}

func (s *Store) Set(key string, resp []byte) error {
	return errors.New("Set() method is not implemented yet")
}

func (s *Store) Delete(key string) error {
	return errors.New("Delete() method is not implemented yet")
}

// Close deletes the storage
func (s *Store) Close() error {
	return nil
}

func (s *Store) setModel(widgetsSchema *dal.Collection) error {

	// register models to this database backend
	widgets := mapper.NewModel(s.backend, widgetsSchema)

	// create the model tables if they don't exist
	if err := widgets.Migrate(); err != nil {
		fmt.Printf("failed to create widget table: %v\n", err)
		return err
	}

	s.mapper = widgets
	return nil
}

func (s *Store) newWidget(widget Widget) {
	// make a new Widget instance, containing the data we want to see
	// the ID field will be populated after creation with the auto-
	// generated UUID.
	newWidget := widget
	/*
		Widget{
			Type:  `foo`,
			Usage: `A fooable widget.`,
		}
	*/

	// insert a widget (ID will be auto-generated because of dal.GenerateUUID)
	if err := Widgets.Create(&newWidget); err != nil {
		fmt.Printf("failed to insert widget: %v\n", err)
		return
	}

	// retrieve the widget using the ID we just got back
	var gotWidget Widget

	if err := Widgets.Get(newWidget.ID, &gotWidget); err != nil {
		fmt.Printf("failed to retrieve widget: %v\n", err)
		return
	}

	fmt.Printf("Got Widget: %#+v", gotWidget)

	// delete the widget
	if err := Widgets.Delete(newWidget.ID); err != nil {
		fmt.Printf("failed to delete widget: %v\n", err)
		return
	}
}
