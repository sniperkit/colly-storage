package compress

const (
	Unknown    Format = iota // unknown format
	GZip                     // Gzip compression format
	BZip2                    // Bzip2 compression
	LZ4                      // LZ4 compression
	Tar                      // Tar format; normally used
	Tar1                     // Tar1 magicnum format; normalizes to Tar
	Tar2                     // Tar1 magicnum format; normalizes to Tar
	Zip                      // Zip archive
	ZipEmpty                 // Empty Zip Archive
	ZipSpanned               // Spanned Zip Archive
)

const formatName = "UnknownGZipBZip2LZ4TarTar1Tar2ZipEmpty ZipSpanned Zip"
