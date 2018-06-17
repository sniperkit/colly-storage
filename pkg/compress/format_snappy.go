package compress

import (
	// external
	"github.com/golang/snappy"
)

func Compress(data []byte) ([]byte, error) {
	return snappy.Encode([]byte{}, data), nil
}

func Decompress(data []byte) ([]byte, error) {
	return snappy.Decode([]byte{}, data)
}
