package badgerstorage

const (

	// storageBucketName...
	storageBucketName string = "colly-store"

	// storagePrefixPath...
	storagePrefixPath string = "./shared/data/storage/badger"

	//-- End
)

const (
	// GzipMinSize gzip min size
	GzipMinSize = 1024

	// StoreFormatRaw raw
	StoreFormatRaw = 0

	// StoreFormatRawGzip raw gzip
	StoreFormatRawGzip = 1

	// StoreFormatJSON json
	StoreFormatJSON = 10

	// StoreFormatJSONGzip json gzip
	StoreFormatJSONGzip = 11

	//-- End
)
