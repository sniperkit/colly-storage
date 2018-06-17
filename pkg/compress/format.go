package compress

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type compressed struct {
	bufBr   *bytes.Buffer
	bufRaw  *bytes.Buffer
	bufGzip *bytes.Buffer
}

type Format int

var formatIndex = [...]uint8{0, 7, 11, 16, 19, 22, 26, 30, 33, 42, 53}

func (i Format) String() string {
	if i < 0 || i >= Format(len(formatIndex)-1) {
		return fmt.Sprintf("Format(%d)", i)
	}
	return formatName[formatIndex[i]:formatIndex[i+1]]
}

// Ext returns the extension for the format. Formats may have more than one
// accepted extension; alternate extensiona are not supported.
func (f Format) Ext() string {
	switch f {
	case GZip:
		return ".gz"
	case BZip2:
		return ".bz2"
	case LZ4:
		return ".lz4"
	case Tar, Tar1, Tar2:
		return ".tar"
	case Zip, ZipEmpty, ZipSpanned:
		return ".zip"
		//case LZW:
		//	return ".Z"
	}
	return "unknown"
}

// ParseFormat takes a string and returns the format or unknown. Any compressed
// tar extensions are returned as the compression format and not tar.
//
// If the passed string starts with a '.', it is removed.
// All strings are lowercased
func ParseFormat(s string) Format {
	if len(s) == 0 {
		return Unknown
	}
	if s[0] == '.' {
		s = s[1:]
	}
	s = strings.ToLower(s)
	switch s {
	case "gzip", "tar.gz", "tgz":
		return GZip
	case "tar":
		return Tar
	case "bz2", "tbz", "tb2", "tbz2", "tar.bz2":
		return BZip2
	case "lz4", "tar.lz4", "tz4":
		return LZ4
	case "zip":
		return Zip
	}
	return Unknown
}

// GetFormat tries to match up the data in the Reader to a supported
// magic number, if a match isn't found, UnsupportedFmt is returned
//
// For zips, this will also match on files with empty zip or spanned zip magic
// numbers.  If you need to distinguich between the various zip formats, use
// something else.
func GetFormat(r io.ReaderAt) (Format, error) {
	// see if the reader contains anything
	b := make([]byte, 1)
	if _, err := r.ReadAt(b, 0); err == io.EOF {
		return Unknown, ErrEmpty
	}

	ok, err := IsLZ4(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return LZ4, nil
	}
	ok, err = IsGZip(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return GZip, nil
	}
	ok, err = IsZip(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return Zip, nil
	}
	ok, err = IsTar(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return Tar, nil
	}
	ok, err = IsBZip2(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return BZip2, nil
	}
	//ok, err = IsLZW(r)
	//if err != nil {
	//	return Unknown, err
	//}
	//if ok {
	//	return LZW, nil
	//}
	return Unknown, ErrUnknown
}
