package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

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
