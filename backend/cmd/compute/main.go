package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	addr         string
	hash         string
	max          int
	workerCount  int
	timeoutInSec int
	dbConn       db.DBStoreClient
}

func main() {
	var c client
	info, _ := debug.ReadBuildInfo()
	for _, bs := range info.Settings {
		if bs.Key == "vcs.revision" {
			c.hash = bs.Value
		}
	}
	flag.IntVar(&c.max, "max", -1, "max number of entries to compute")
	flag.IntVar(&c.workerCount, "w", 10, "number of workers to use")
	flag.IntVar(&c.timeoutInSec, "timeout", 30, "time out in seconds to run each sim for")
	flag.Parse()
	//compute steps
	// 1. ask server for work
	// 2. compute work
	// 3. ask again for more work
	c.addr = os.Getenv("DB_RPC_ADDR")

	if c.addr == "" {
		log.Fatal("Invalid address for rpc service")
	}

	log.Printf("starting compute run, hash: %v, addr %v", c.hash, c.addr)
	err := c.run()
	if err != nil {
		log.Fatal(err)
	}

}

func (c *client) run() error {
	conn, err := grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c.dbConn = db.NewDBStoreClient(conn)
	for {
		work, err := c.getWork()
		if err != nil {
			return err
		}
		if len(work) == 0 {
			log.Println("all done")
			return nil
		}
		err = c.processBatch(work)
		if err != nil {
			//this is only if we have unexpected error like rpc failed?
			return err
		}
		//stop if we're done up to max
		if c.max == 0 {
			return nil
		}
	}
}

func (c *client) getWork() ([]*db.ComputeWork, error) {
	resp, err := c.dbConn.GetWork(context.Background(), &db.GetWorkRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetData(), nil
}

// return number of work item completed
func (c *client) processBatch(w []*db.ComputeWork) error {
	for _, v := range w {
		res, err := c.processWork(v)
		if err != nil {
			err = c.postError(v, err.Error())
		} else {
			err = c.postResult(v, res)
		}
		if err != nil {
			return err
		}

		//don't do too much work
		if c.max != -1 {
			c.max--
			if c.max == 0 {
				return nil
			}
		}
	}
	return nil
}

func (c *client) processWork(w *db.ComputeWork) (*model.SimulationResult, error) {
	start := time.Now()
	//compute work??
	log.Printf("got work %v; starting compute", w.Id)
	// compute result
	simcfg, err := simulator.Parse(w.Config)
	if err != nil {
		log.Printf("could not parse config for id %v: %v\n", w.Id, err)
		//TODO: we should post something here??
		return nil, err
	}
	simcfg.Settings.Iterations = int(w.Iterations)
	simcfg.Settings.NumberOfWorkers = 30

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeoutInSec)*time.Second)
	defer cancel()

	result, err := simulator.RunWithConfig(w.Config, simcfg, simulator.Options{}, time.Now(), ctx)
	if err != nil {
		log.Printf("error running sim %v: %v\n", w.Id, err)
		return nil, err
	}
	elapsed := time.Since(start)
	log.Printf("Work %v took %s", w.Id, elapsed)

	return result, nil
}

func (c *client) postError(w *db.ComputeWork, reason string) error {
	_, err := c.dbConn.RejectWork(context.Background(), &db.RejectWorkRequest{
		Id:     w.Id,
		Reason: reason,
		Hash:   c.hash,
	})
	return err
}

func (c *client) postResult(w *db.ComputeWork, res *model.SimulationResult) error {
	_, err := c.dbConn.CompleteWork(context.Background(), &db.CompleteWorkRequest{
		Id:     w.Id,
		Result: res,
	})
	return err
}
