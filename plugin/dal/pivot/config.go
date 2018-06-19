package dal_pivot

import (
	"github.com/imdario/mergo"
)

// Config is...
type Config struct {
	Identifier string                 `json:"identifier" yaml:"identifier" toml:"identifier" xml:"identifier" ini:"identifier" config:"store.dal.identifier"`
	Shorthand  string                 `json:"shorthand" yaml:"shorthand" toml:"shorthand" xml:"shorthand" ini:"shorthand" config:"store.dal.shorthand"`
	DSN        string                 `json:"dsn" yaml:"dsn" toml:"dsn" xml:"dsn" ini:"dsn" config:"store.dal.dsn"`
	Scheme     string                 `default:"sqlite" json:"scheme" yaml:"scheme" toml:"scheme" xml:"scheme" ini:"scheme" config:"store.dal.scheme"`
	Host       string                 `default:"./shared/storage/dal/pivot/colly.db" json:"host" yaml:"host" toml:"host" xml:"host" ini:"host" config:"store.dal.host"`
	Dataset    string                 `default:"colly-data" json:"dataset" yaml:"dataset" toml:"dataset" xml:"dataset" ini:"dataset" config:"store.dal.dataset"`
	Options    map[string]interface{} `json:"options" yaml:"options" toml:"options" xml:"options" ini:"options" config:"store.dal.options"`
	Sanitize   bool                   `json:"sanitize" yaml:"sanitize" config:"store.dal.sanitize" config:"store.dal.sanitize"`
	Debug      bool                   `json:"debug" yaml:"debug" config:"store.dal.debug" config:"store.dal.debug"`
	Verbose    bool                   `json:"verbose" yaml:"verbose" config:"store.dal.verbose" config:"store.dal.verbose"`
	done       chan struct{}          `json:"-" yaml:"-" toml:"-" xml:"-" config:"-" config:"-"`
}

type backendConfig struct {
	Engine string `json:"provider" yaml:"provider" config:"store.dal.provider"`
}

// DefaultConfig returns the default configuration for this serializer
func DefaultConfig() Config {
	return Config{
		Sanitize: false,
		DSN:      DefaultBackendDSN,
	}
}

// ConfigName ...
func (Config) ConfigName() string {
	return "DAL"
}

// Merge merges the default with the given config and returns the result
func (c Config) Merge(cfg []Config) (config Config) {
	if len(cfg) > 0 {
		config = cfg[0]
		mergo.Merge(&config, c)
	} else {
		_default := c
		config = _default
	}
	return
}

// MergeSingle merges the default with the given config and returns the result
func (c Config) MergeSingle(cfg Config) (config Config) {
	config = cfg
	mergo.Merge(&config, c)
	return
}

/*
// storageConfig is...
type storageConfig struct {
	Provider       string        `json:"provider" yaml:"provider" config:"store.dal.provider"`
	MaxConnections int           `json:"max_connections" yaml:"max_connections" config:"store.dal.max_connections" default:"0"`
	BucketName     string        `json:"bucket_name" yaml:"bucket_name" config:"store.dal.provider"`
	StoragePath    string        `json:"storage_path" yaml:"storage_path" config:"store.dal.storage_path"`
	EnableGzip     bool          `json:"enable_gzip" yaml:"enable_gzip" config:"store.dal.enable_gzip"`
	ReadOnly       bool          `json:"read_only" yaml:"read_only" config:"store.dal.read_only"`
	StrictMode     bool          `json:"strict_mode" yaml:"strict_mode" config:"store.dal.strict_mode"`
	NoSync         bool          `json:"no_sync" yaml:"no_sync" config:"store.dal.no_sync"`
	NoFreelistSync bool          `json:"no_freelist_sync" yaml:"no_freelist_sync" config:"store.dal.no_freelist_sync"`
	NoGrowSync     bool          `json:"no_grow_sync" yaml:"no_grow_sync" config:"store.dal.no_grow_sync"`
	MaxBatchSize   bool          `json:"max_batch_size" yaml:"max_batch_size" config:"store.dal.max_batch_size"`
	MaxBatchDelay  bool          `json:"max_batch_delay" yaml:"max_batch_delay" config:"store.dal.max_batch_delay"`
	AllocSize      bool          `json:"alloc_size" yaml:"alloc_size" config:"store.dal.provialloc_sizeder"`
	done           chan struct{} `json:"-" yaml:"-" toml:"-" xml:"-" config:"-"`
}

// Wait ...
func (c storageConfig) Wait() {
	<-c.done
}
*/
