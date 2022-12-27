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
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulator"
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

func main() {

	var data []entry

	err := getJson("https://handleronedev.gcsim.app/db_simulations", &data)
	if err != nil {
		panic(err)
	}

	// log.Println(data)

	client, err := db.NewClient(db.ClientCfg{
		Addr: "localhost:3000",
	})
	if err != nil {
		panic(err)
	}

	//loop through, rerun, and store in mongo?
	if len(data) == 0 {
		log.Println("no data; exiting")
	}

	for _, v := range data[:10] {
		start := time.Now()
		log.Printf("Recomputing %v\n", v.Key)
		e, err := parseAndComputeEntry(v)
		if err != nil {
			log.Printf("Skipping db entry: %v; error encountered: %v\n", v.Key, err)
			continue
		}
		key, err := client.Create(context.TODO(), e)

		if err != nil {
			log.Println("insert into db failed! panicking")
			panic(err)
		}

		elapsed := time.Since(start)
		log.Printf("%v completed in %s; new entry with key %v", v.Key, elapsed, key)
	}

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

func parseAndComputeEntry(e entry) (*model.DBEntry, error) {
	z, err := base64.StdEncoding.DecodeString(e.ConfigHash)
	if err != nil {
		return nil, err
	}
	r, err := zlib.NewReader(bytes.NewReader(z))
	if err != nil {
		return nil, err
	}
	cfg, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	simcfg, err := simulator.Parse(string(cfg))
	if err != nil {
		return nil, err
	}

	simcfg.Settings.Iterations = iters
	simcfg.Settings.NumberOfWorkers = workers

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	summary, err := simulator.RunWithConfig(string(cfg), simcfg, simulator.Options{}, time.Now(), ctx)
	if err != nil {
		return nil, err
	}

	return summary.ToPBDBEntry(), nil
}
