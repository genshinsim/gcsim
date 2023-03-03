package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/genshinsim/gcsim/backend/pkg/services/submission"
	"github.com/genshinsim/gcsim/pkg/model"
)

type entry struct {
	ID         int    `json:"db_id"`
	Key        string `json:"simulation_key"`
	Hash       string `json:"git_hash"`
	ConfigHash string `json:"config_hash"`
	Desc       string `json:"sim_description"`
}

const iters = 100
const workers = 30

var client *submission.Client

func main() {

	var data []entry

	err := getJson("https://handleronedev.gcsim.app/db_simulations", &data)
	if err != nil {
		panic(err)
	}

	// log.Println(data)

	client, err = submission.NewClient("192.168.100.47:8083")
	if err != nil {
		panic(err)
	}

	//loop through, rerun, and store in mongo?
	if len(data) == 0 {
		log.Println("no data; exiting")
	}

	for i, v := range data[:10] {
		err := addSubmission(v, i)
		if err != nil {
			log.Printf("Skipping db entry: %v; error encountered: %v\n", v.Key, err)
		}
	}

}

func addSubmission(v entry, i int) error {
	z, err := base64.StdEncoding.DecodeString(v.ConfigHash)
	if err != nil {
		return err
	}
	r, err := zlib.NewReader(bytes.NewReader(z))
	if err != nil {
		return err
	}
	cfg, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	sub := &model.Submission{
		Config:      string(cfg),
		Submitter:   "gcsim#0000",
		Description: v.Desc,
	}
	key, err := client.Submit(context.TODO(), sub)

	if err != nil {
		return err
	}

	log.Printf("%v submitted; new entry with key %v", v.Key, key)

	return nil
}

var myClient = &http.Client{Timeout: 120 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
