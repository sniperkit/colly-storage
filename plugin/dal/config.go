package dal

import (
	"time"

	"github.com/imdario/mergo"
)

// Config is...
type Config struct {
	Engine         string        `json:"provider" yaml:"provider" config:"store.http.provider"`
	DatabaseName   *string       `json:"bucket_name" yaml:"bucket_name" config:"store.http.provider"`
	BucketName     *string       `json:"bucket_name" yaml:"bucket_name" config:"store.http.provider"`
	StoragePath    *string       `json:"storage_path" yaml:"storage_path" config:"store.http.storage_path"`
	MaxConnections int           `json:"max_connections" yaml:"max_connections" config:"store.http.max_connections" default:"0"`
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

// Wait ...
func (c Config) Wait() {
	<-c.done
}
