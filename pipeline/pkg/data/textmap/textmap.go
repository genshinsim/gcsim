package textmap

import (
	"encoding/json"
	"fmt"
	"os"
)

type DataSource struct {
	TextMap map[int64]string
}

func NewTextMapSource(paths []string) (*DataSource, error) {
	var res DataSource

	for _, path := range paths {
		d, err := os.ReadFile(path)
		// if error, try next path, print warning
		if err != nil {
			continue
		}
		err = json.Unmarshal(d, &res.TextMap)
		if err != nil {
			return nil, err
		}
	}
	if res.TextMap == nil {
		return nil, fmt.Errorf("error reading textmap files: no valid files found")
	}

	return &res, nil
}

func (d *DataSource) Get(id int64) (string, error) {
	val, ok := d.TextMap[id]
	if !ok {
		return "", fmt.Errorf("could not find data for id %v", id)
	}
	return val, nil
}
