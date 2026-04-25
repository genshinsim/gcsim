package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type dbGetOpt struct {
	Query   map[string]any `json:"query,omitempty"`
	Sort    map[string]any `json:"sort,omitempty"`
	Project map[string]any `json:"project,omitempty"`
	Skip    int64          `json:"skip,omitempty"`
	Limit   int64          `json:"limit,omitempty"`
}

type dbData struct {
	Data []dbEntry `json:"data"`
}

type dbEntry struct {
	Id     string `json:"_id"`
	Config string `json:"config"`
}

func getDBEntries() ([]dbEntry, error) {
	q := dbGetOpt{
		Limit: 100,
	}

	var page int64
	var res []dbEntry

	for {
		q.Skip = 100 * page
		jsonStr, err := json.Marshal(q)
		if err != nil {
			return nil, err
		}
		url := fmt.Sprintf("https://simpact.app/api/db?q=%v", url.QueryEscape(string(jsonStr)))

		var d dbData
		err = getJSON(url, &d)
		if err != nil {
			return nil, err
		}
		if len(d.Data) == 0 {
			break
		}
		res = append(res, d.Data...)

		page++
	}

	return res, nil
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, target any) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
