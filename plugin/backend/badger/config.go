package badgerstorage

import (
	"github.com/imdario/mergo"
)

// Config is the configuration for this serializer
type Config struct {
	Debug bool

	Sanitize bool

	Compress bool

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

	// 4. Flags for testing purposes
	// ------------------------------
	DoNotCompact bool // Stops LSM tree from compactions.

}

// DefaultConfig returns the default configuration for this serializer
func DefaultConfig() Config {
	return Config{
		Sanitize: false,
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
