package dal

import (
	"fmt"
	"sync"
	"time"

	"github.com/ghetzel/pivot"
	"github.com/ghetzel/pivot/backends"
	"github.com/ghetzel/pivot/dal"
	"github.com/ghetzel/pivot/mapper"
)

var (
	DefaultBackendDSN string = fmt.Printf("sqlite:///%s/%s.%s", StoragePrefixPath, StorageDatabaseName, DefaultStorageFileExtension)
)
