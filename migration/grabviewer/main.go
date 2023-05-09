package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type entry struct {
	Key        string    `json:"simulation_key" bson:"_id"`
	Metadata   string    `json:"metadata" bson:"metadata"`
	Viewer     string    `json:"viewer_file" bson:"viewer_file"`
	Permanent  bool      `json:"is_permanent" bson:"is_permanent"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

const limit = 25
const max = 10000

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
)

func main() {
	err := mongoConnect("mongodb://192.168.100.102:2700")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(exitCodeInterrupt)
	}()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitCodeErr)
	}

}

func run(ctx context.Context) error {
	//find last inserted
	last, err := mongoFindLatest()
	if err != nil {
		panic(err)
	}
	log.Printf("last entry: %v - %v\n", last.Key, last.CreateTime)
	count := 0

	format := "2006-01-02T15:04:05.999999999"

	for {
		select {
		case <-ctx.Done():
			log.Println("interrupt signal received")
			return nil
		default:
			if count > max {
				log.Println("limit reached")
				return nil
			}

			var data []entry

			q := fmt.Sprintf(`https://handleronedev.gcsim.app/simulations?order=create_time.asc&create_time=gte.%v&limit=%v`, last.CreateTime.Format(format), limit)
			log.Printf("starting batch; query: %v\n", q)

			err = getData(context.TODO(), q, &data)
			if err != nil {
				panic(err)
			}

			log.Println("query ok")

			if len(data) == 0 {
				log.Println("no more data!")
				return nil
			}

			if len(data) == 1 && data[0].Key == last.Key {
				log.Println("no more data!")
				return nil
			}

			for _, v := range data {
				if v.Key == last.Key {
					// log.Println("skipping last")
					continue
				}
				// log.Println(v.Key, v.CreateTime)
				err = mongoInsert(context.TODO(), v)
				if err != nil {
					return err
				}
			}
			last = data[len(data)-1]
			count += len(data)
			log.Printf("inserted %v, count is now %v\n", len(data), count)
		}

	}
}

func getData(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

var col *mongo.Collection

func mongoConnect(addr string) error {
	credential := options.Credential{
		Username: "root",
		Password: "root-password",
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(addr).SetAuth(credential))
	if err != nil {
		return err
	}

	//check connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Println("pig failed", err)
		return err
	}

	col = client.Database("store").Collection("data")
	return nil
}

func mongoInsert(ctx context.Context, d entry) error {
	res, err := col.InsertOne(ctx, d)
	if err != nil {
		return err
	}
	log.Println("inserted with id: ", res.InsertedID)
	return nil
}

func mongoFindLatest() (entry, error) {
	var result entry
	res := col.FindOne(context.TODO(), bson.M{}, &options.FindOneOptions{
		Sort: bson.M{
			"create_time": -1,
		},
	})
	err := res.Err()
	if err != nil {
		return result, err
	}
	err = res.Decode(&result)
	if err != nil {
		log.Println("error decoding")
		return result, err
	}
	return result, nil
}
