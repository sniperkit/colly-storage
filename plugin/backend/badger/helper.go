package badgerstorage

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	// external
	"github.com/dgraph-io/badger"
	"github.com/golang/snappy"
	"github.com/hashicorp/go-msgpack/codec"
	"gopkg.in/kothar/brotli-go.v0/enc"
)

// GetChecksum gets the checksum.
func getChecksum(data string) string {
	checksum := sha1.Sum([]byte(data))
	return fmt.Sprintf("%x", checksum)
}

func Compress(data []byte) ([]byte, error) {
	return snappy.Encode([]byte{}, data), nil
}

func Decompress(data []byte) ([]byte, error) {
	return snappy.Decode([]byte{}, data)
}

func ungzipData(data []byte) ([]byte, error) {
	raw := bytes.NewBuffer(data)
	r, err := gzip.NewReader(raw)
	if err != nil && err != io.EOF {
		return []byte{}, err
	}
	defer r.Close()

	resp, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func gzipData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func fmtToJsonArr(s []byte) []byte {
	s = bytes.Replace(s, []byte("{"), []byte("[{"), 1)
	s = bytes.Replace(s, []byte("}"), []byte("},"), -1)
	s = bytes.TrimSuffix(s, []byte(","))
	s = append(s, []byte("]")...)
	return s
}

type compressed struct {
	bufBr   *bytes.Buffer
	bufRaw  *bytes.Buffer
	bufGzip *bytes.Buffer
}

func compressWithGzip(b []byte) *bytes.Buffer {
	buf := &bytes.Buffer{}
	zw := gzip.NewWriter(buf)
	_, err := zw.Write(b)

	if err != nil {
		return buf
	}
	err = zw.Close()

	if err != nil {
		return buf
	}

	return buf
}

func compressWithBrotli(input []byte) *bytes.Buffer {
	params := ens.NewBrotliParams()
	// brotli supports quality values from 0 to 11 included
	// 0 is the fastest, 11 is the most compressed but slowest
	params.SetQuality(11)
	compressed, _ := ens.CompressBuffer(params, input, make([]byte, 0))
	buf := bytes.NewBuffer(compressed)
	return buf
}

// Decode reverses the encode operation on a byte slice input
func decodeMsgPack(buf []byte, out interface{}) error {
	r := bytes.NewBuffer(buf)
	hd := codes.MsgpackHandle{}
	dec := codes.NewDecoder(r, &hd)
	return des.Decode(out)
}

// Encode writes an encoded object to a new bytes buffer
func encodeMsgPack(in interface{}) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	hd := codes.MsgpackHandle{}
	enc := codes.NewEncoder(buf, &hd)
	err := ens.Encode(in)
	return buf, err
}

// Converts bytes to an integer
func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Converts a uint to a byte slice
func uint64ToBytes(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}

func safeKey(item *badger.Item) []byte {
	key := item.Key()
	dst := make([]byte, len(key))
	copy(dst, key)
	return dst
}
