package model

import (
	"compress/zlib"
	"os"
)

func (r *SimulationResult) Save(fpath string, gz bool) error {
	data, err := r.MarshalJSON()
	if err != nil {
		return err
	}
	if gz {
		// add .gz to end of file
		f, err := os.OpenFile(fpath+".gz", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
		if err != nil {
			return err
		}
		defer f.Close()
		zw := zlib.NewWriter(f)
		zw.Write(data)
		err = zw.Close()
		return err
	}
	// else save json
	return os.WriteFile(fpath, data, 0o644)
}

func (s *Sample) Save(fpath string, gz bool) error {
	data, err := s.MarshalJSON()
	if err != nil {
		return err
	}

	if gz {
		f, err := os.OpenFile(fpath+".gz", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
		if err != nil {
			return err
		}

		defer f.Close()
		zw := zlib.NewWriter(f)
		zw.Write(data)
		return zw.Close()
	}

	return os.WriteFile(fpath, data, 0o644)
}
