package model

import (
	"compress/zlib"
	"os"
)

func (s *Sample) Save(fpath string, gz bool) error {
	data, err := s.MarshalJSON()
	if err != nil {
		return err
	}
	return WriteJSON(fpath, data, gz)
}

// WriteJSON writes data to fpath, or zlib-compressed to fpath.gz.
func WriteJSON(fpath string, data []byte, gz bool) error {
	if !gz {
		return os.WriteFile(fpath, data, 0o644)
	}
	f, err := os.OpenFile(fpath+".gz", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zlib.NewWriter(f)
	if _, err := zw.Write(data); err != nil {
		return err
	}
	return zw.Close()
}
