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

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type entry struct {
	ID         int    `json:"db_id"`
	Key        string `json:"simulation_key"`
	Hash       string `json:"git_hash"`
	ConfigHash string `json:"config_hash"`
	Desc       string `json:"sim_description"`
}

var client db.DBStoreClient

func main() {

	var data []entry

	err := getJson("https://handleronedev.gcsim.app/db_simulations", &data)
	if err != nil {
		panic(err)
	}

	// log.Println(data)
	conn, err := grpc.Dial("192.168.100.102:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client = db.NewDBStoreClient(conn)

	//loop through, rerun, and store in mongo?
	if len(data) == 0 {
		log.Println("no data; exiting")
	}

	for i, v := range data[:50] {
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
	resp, err := client.Submit(context.TODO(), &db.SubmitRequest{
		Config:      string(cfg),
		Submitter:   "0000",
		Description: v.Desc,
	})
	if err != nil {
		return err
	}

	log.Printf("%v submitted; new entry with key %v", v.Key, resp.GetId())

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
