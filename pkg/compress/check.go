package compress

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"

	compressible "github.com/sniperkit/compressible/pkg"
)

// Magic numbers for compression and archive formats
var (
	magicnumGZip       = []byte{0x1f, 0x8b}
	magicnumBZip2      = []byte{0x42, 0x5a, 0x68}
	magicnumLZ4        = []byte{0x18, 0x4d, 0x22, 0x04}
	magicnumTar1       = []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30} // offset: 257
	magicnumTar2       = []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x20, 0x00} // offset: 257
	magicnumZip        = []byte{0x50, 0x4b, 0x03, 0x04}
	magicnumZipEmpty   = []byte{0x50, 0x4b, 0x05, 0x06}
	magicnumZipSpanned = []byte{0x50, 0x4b, 0x07, 0x08}
	//magicnumLZW        = []byte{0x1F, 0x9d}
)

// IsBZip2 checks to see if the received reader's contents are in bzip2 format
// by checking the magic numbers.
func IsBZip2(r io.ReaderAt) (bool, error) {
	h := make([]byte, 3)
	// Read the first 3 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var hb [3]byte
	// check for bzip2
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &hb)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched bzip2's magic number: %s", err)
	}
	var cb [3]byte
	cbuf := bytes.NewBuffer(magicnumBZip2)
	err = binary.Read(cbuf, binary.BigEndian, &cb)
	if err != nil {
		return false, fmt.Errorf("error while converting bzip2 magic number for comparison: %s", err)
	}
	if hb == cb {
		return true, nil
	}
	return false, nil
}

// IsGZip checks to see if the received reader's contents are in gzip format
// by checking the magic numbers.
func IsGZip(r io.ReaderAt) (bool, error) {
	h := make([]byte, 2)
	// Read the first 2 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h16 uint16
	// check for gzip
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.BigEndian, &h16)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched bzip2's magic number: %s", err)
	}
	var c16 uint16
	cbuf := bytes.NewBuffer(magicnumGZip)
	err = binary.Read(cbuf, binary.BigEndian, &c16)
	if err != nil {
		return false, fmt.Errorf("error while converting bzip2 magic number for comparison: %s", err)
	}
	if h16 == c16 {
		return true, nil
	}
	return false, nil
}

// IsLZ4 checks to see if the received reader's contents are in LZ4 foramt by
// checking the magic numbers.
func IsLZ4(r io.ReaderAt) (bool, error) {
	h := make([]byte, 4)
	// Read the first 4 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h32 uint32
	// check for lz4
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &h32)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched LZ4's magic number: %s", err)
	}
	var c32 uint32
	cbuf := bytes.NewBuffer(magicnumLZ4)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting LZ4 magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	return false, nil
}

// IsLZW checks to see if the received reader's contents are in LZ4 format by
// checking the magic numbers.
//
// TODO: unsupported until I have a better understanding of how to handle LZW
/*
func IsLZW(r io.ReaderAt) (bool, error) {
	h := make([]byte, 2)
	// Reat the first 8 bytes since that's where most magic numbers are
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h16 uint16
	// check for lzw
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &h16)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched LZW's magic number: %s", err)
	}
	var c16 uint16
	cbuf := bytes.NewBuffer(magicnumLZW)
	err = binary.Read(cbuf, binary.BigEndian, &c16)
	if err != nil {
		return false, fmt.Errorf("error while converting LZW magic number for comparison: %s", err)
	}
	if h16 == c16 {
		return true, nil
	}
	return false, nil
}
*/

// IsTar checks to see if the received reader's contents are in the tar format
// by checking the magic numbers. This evaluates using both tar1 and tar2 magic
// numbers.
func IsTar(r io.ReaderAt) (bool, error) {
	h := make([]byte, 8)
	// Read the first 8 bytes at offset 257
	_, err := r.ReadAt(h, 257)
	if err != nil {
		return false, err
	}
	var h64 uint64
	// check for Zip
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.BigEndian, &h64)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched tar's magic number: %s", err)
	}
	var c64 uint64
	cbuf := bytes.NewBuffer(magicnumTar1)
	err = binary.Read(cbuf, binary.BigEndian, &c64)
	if err != nil {
		return false, fmt.Errorf("error while converting the tar magic number for comparison: %s", err)
	}
	if h64 == c64 {
		return true, nil
	}
	cbuf = bytes.NewBuffer(magicnumTar2)
	err = binary.Read(cbuf, binary.BigEndian, &c64)
	if err != nil {
		return false, fmt.Errorf("error while converting the empty tar magic number for comparison: %s", err)
	}
	if h64 == c64 {
		return true, nil
	}
	return false, nil
}

// IsZip checks to see if the received reader's contents are in the zip format
// by checking the magic numbers. This will match on zip, empty zip and spanned
// zip magic numbers. If you need to distinguish between those, use something
// else.
func IsZip(r io.ReaderAt) (bool, error) {
	h := make([]byte, 4)
	// Read the first 4 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h32 uint32
	// check for Zip
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.BigEndian, &h32)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched zip's magic number: %s", err)
	}
	var c32 uint32
	cbuf := bytes.NewBuffer(magicnumZip)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting the zip magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	cbuf = bytes.NewBuffer(magicnumZipEmpty)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting the empty zip magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	cbuf = bytes.NewBuffer(magicnumZipSpanned)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting the spanned zip magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	return false, nil
}

func IsCompressible(ct string) bool {
	return compressible.Is(ct)
}

func IsCompressibleWithThreshold(ct string, wt compressible.WithThreshold) bool {
	return wt.Compressible(ct, wt)
}

// GetChecksum gets the checksum.
func getChecksum(data string) string {
	checksum := sha1.Sum([]byte(data))
	return fmt.Sprintf("%x", checksum)
}
