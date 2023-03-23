package data

import (
	"encoding/json"
	"os"
)

func loadAvatarExcel(path string) ([]AvatarExcel, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var res []AvatarExcel
	err = json.Unmarshal(d, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
