package textmap

import (
	"encoding/json"
	"fmt"
	"os"
)

type DataSource struct {
	TextMap map[int64]string
}

func NewTextMapSource(path string) (*DataSource, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var res DataSource
	err = json.Unmarshal(d, &res.TextMap)
	if err != nil {
		return nil, err
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
