package result

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

func (s *Summary) Save(fpath string, gz bool) error {
	s.Text = s.PrettyPrint()
	s.V2 = true //force v2
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if gz {
		//add .gz to end of file
		f, err := os.OpenFile(fpath+".gz", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		zw := gzip.NewWriter(f)
		zw.Write(data)
		err = zw.Close()
		return err
	}
	//else save json
	return os.WriteFile(fpath, data, 0644)
}
