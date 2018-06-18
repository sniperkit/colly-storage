package dal_pivot

import (
	"errors"
	"fmt"
	"sync"

	// external
	"github.com/ghetzel/pivot"
	"github.com/ghetzel/pivot/backends"
	"github.com/ghetzel/pivot/dal"
	"github.com/ghetzel/pivot/mapper"

	// internal
	helper "github.com/sniperkit/colly-storage/pkg/helper"
)

var (

	// sqlite:///./test.db
	DefaultBackendDSN string = fmt.Sprintf(`sqlite:///%s/%s%s`, DefaultStoragePrefixPath, DefaultStorageDatabaseName, DefaultStorageFileExtension)
)

type Store struct {
	lock    *sync.RWMutex
	backend backends.Backend
	mapper  mapper.Mapper
	schema  []*dal.Collection
}

func NewDataAbstractionLayer(config *Config) (*Store, error) {
	s := &Store{
		lock: &sync.RWMutex{},
	}

	if config.DSN == "" {
		config.DSN = DefaultBackendDSN
	}

	if config.PrefixPath == nil {
		config.PrefixPath = *DefaultStoragePrefixPath
	}

	if config.PrefixPath != nil {
		if err := helper.EnsurePathExists(*config.PrefixPath); err != nil {
			return nil, err
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

	return s, nil
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
