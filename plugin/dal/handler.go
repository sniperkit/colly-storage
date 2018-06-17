package dal

import (
	"sync"

	"github.com/ghetzel/pivot"
	"github.com/ghetzel/pivot/backends"
	"github.com/ghetzel/pivot/dal"
)

type Storage struct {
	lock   *sync.RWMutex
	store  backends.Backend
	schema []*dal.Collection
}

type Config struct {
	DSN      string
	ReadOnly bool
	Compress bool
	Debug    bool
	Stats    bool
}

func NewDataAbstractionLayer(config *Config) (*Store, error) {
	s := &Store{
		lock: &sync.RWMutex{},
	}

	if config.DSN == "" {
		config.DSN = DefaultBackendDSN
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

// Init initializes Store
func (s *Store) Init() error {
	return errors.New("Init method is not implemented yet...")
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

// Debug
func (s *Store) Debug(action string) error {
	return errors.New("Debug() method is not implemented yet")
}

func (c *Store) Action(name string, args ...interface{}) (map[string]*interface{}, error) {
	return nil, errors.New("Action() method is not implemented yet")
}

// Ping check if the storage is available...
func (s *Store) Ping() error {
	return errors.New("Ping() method is not implemented yet")
}

// Clear truncate all key/values stored...
func (s *Store) Clear() error {
	return errors.New("Debug is not implemented yet")
}

// Close deletes the storage
func (s *Store) Close() error {
	return nil
}

func (s *Store) setModel(widgetsSchema *dal.Collection) (mapper.Mapper, error) {
	// register models to this database backend
	widgets = mapper.NewModel(s.backend, widgetsSchema)
	// create the model tables if they don't exist
	if err := widgets.Migrate(); err != nil {
		fmt.Printf("failed to create widget table: %v\n", err)
		return nil, error
	}
	return widgets, nil
}

/*
func (s *Store) setModel() {
	// make a new Widget instance, containing the data we want to see
	// the ID field will be populated after creation with the auto-
	// generated UUID.
	newWidget := Widget{
		Type:  `foo`,
		Usage: `A fooable widget.`,
	}

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
*/
