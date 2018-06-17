package badgerstorage

import (
	"time"

	"github.com/imdario/mergo"
)

// Config is the configuration for this serializer
type Config struct {
	// 1. Mandatory flags
	// -------------------
	// Directory to store the data in. Should exist and be writable.
	StoragePath string // Dir

	// Directory to store the value log in. Can be the same as Dir. Should
	// exist and be writable.
	ValueDir string

	// 2. Frequently modified flags
	// -----------------------------
	// Sync all writes to disk. Setting this to true would slow down data
	// loading significantly.
	SyncWrites bool

	// 3. Flags that user might want to review
	// ----------------------------------------
	// The following affect all levels of LSM tree.
	MaxTableSize        int64 // Each table (or file) is at most this size.
	LevelSizeMultiplier int   // Equals SizeOf(Li+1)/SizeOf(Li).
	MaxLevels           int   // Maximum number of levels of compaction.

	// If value size >= this threshold, only store value offsets in tree.
	ValueThreshold int

	// Maximum number of tables to keep in memory, before stalling.
	NumMemtables int

	// The following affect how we handle LSM tree L0.
	// Maximum number of Level 0 tables before we start compacting.
	NumLevelZeroTables int

	// If we hit this number of Level 0 tables, we will stall until L0 is
	// compacted away.
	NumLevelZeroTablesStall int

	// Maximum total size for L1.
	LevelOneSize int64

	// Size of single value log file.
	ValueLogFileSize int64

	// Number of compaction workers to run concurrently.
	NumCompactors int

	// 4. Flags for dev purposes
	// ------------------------------
	//
	UseTTL bool

	//
	TTL time.Duration

	// 5. Flags for testing purposes
	// ------------------------------
	DoNotCompact bool // Stops LSM tree from compactions.

	// Debug sets...
	Debug bool

	// Sanitize sets...
	Sanitize bool

	// Compress sets...
	Compress bool
}

// DefaultConfig returns the default configuration for this serializer
func DefaultConfig() Config {
	return Config{
		Sanitize:   false,
		SyncWrites: true,
	}
}

// ConfigName ...
func (Config) ConfigName() string {
	return "BadgerKV"
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

// Wait ...
func (c storageConfig) Wait() {
	<-c.done
}
*/
