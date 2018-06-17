package compress

import (
	"bytes"

	// external
	"gopkg.in/kothar/brotli-go.v0/enc"
)

func compressWithBrotli(input []byte) *bytes.Buffer {
	params := enc.NewBrotliParams()
	// brotli supports quality values from 0 to 11 included
	// 0 is the fastest, 11 is the most compressed but slowest
	params.SetQuality(11)
	compressed, _ := enc.CompressBuffer(params, input, make([]byte, 0))
	buf := bytes.NewBuffer(compressed)
	return buf
}
