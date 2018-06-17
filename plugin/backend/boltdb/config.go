package boltdbstorage

import (
	"github.com/imdario/mergo"
)

// Config is the configuration for this serializer
type Config struct {
	Sanitize    bool
	BucketName  string
	StoragePath string
	Debug       bool
}

// storageConfig is...
type storageConfig struct {
	Provider       string        `json:"provider" yaml:"provider" config:"store.http.provider"`
	MaxConnections int           `json:"max_connections" yaml:"max_connections" config:"store.http.max_connections" default:"0"`
	BucketName     string        `json:"bucket_name" yaml:"bucket_name" config:"store.http.provider"`
	StoragePath    string        `json:"storage_path" yaml:"storage_path" config:"store.http.storage_path"`
	EnableGzip     bool          `json:"enable_gzip" yaml:"enable_gzip" config:"store.http.enable_gzip"`
	ReadOnly       bool          `json:"read_only" yaml:"read_only" config:"store.http.read_only"`
	StrictMode     bool          `json:"strict_mode" yaml:"strict_mode" config:"store.http.strict_mode"`
	NoSync         bool          `json:"no_sync" yaml:"no_sync" config:"store.http.no_sync"`
	NoFreelistSync bool          `json:"no_freelist_sync" yaml:"no_freelist_sync" config:"store.http.no_freelist_sync"`
	NoGrowSync     bool          `json:"no_grow_sync" yaml:"no_grow_sync" config:"store.http.no_grow_sync"`
	MaxBatchSize   bool          `json:"max_batch_size" yaml:"max_batch_size" config:"store.http.max_batch_size"`
	MaxBatchDelay  bool          `json:"max_batch_delay" yaml:"max_batch_delay" config:"store.http.max_batch_delay"`
	AllocSize      bool          `json:"alloc_size" yaml:"alloc_size" config:"store.http.provialloc_sizeder"`
	done           chan struct{} `json:"-" yaml:"-" toml:"-" xml:"-" config:"-"`
}

// DefaultConfig returns the default configuration for this serializer
func DefaultConfig() Config {
	return Config{
		Sanitize:    false,
		StoragePath: defaultStorageFile,
		BucketName:  storageBucketName,
	}
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
