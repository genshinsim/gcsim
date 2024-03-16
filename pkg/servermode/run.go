package servermode

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
)

func errorRecover(r interface{}) error {
	var err error
	switch x := r.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = errors.New("unknown error")
	}
	return err
}

// isRunning checks if an id is running
func (s *Server) isRunning(id string) bool {
	//WARNING: READ ONLY HERE NOT SAFE FOR WRITE
	_, ok := s.pool[id]
	return ok
}

func flush(aggregators []agg.Aggregator) *model.SimulationStatistics {
	stats := &model.SimulationStatistics{}
	for _, a := range aggregators {
		a.Flush(stats)
	}
	return stats
}

func parse(cfg string) (*info.ActionList, ast.Node, error) {
	parser := ast.New(cfg)
	simcfg, gcsl, err := parser.Parse()
	if err != nil {
		return &info.ActionList{}, nil, err
	}

	// check other errors as well
	if len(simcfg.Errors) != 0 {
		fmt.Println("The config has the following errors: ")
		errMsgs := ""
		for _, v := range simcfg.Errors {
			errMsg := fmt.Sprintf("\t%v\n", v)
			fmt.Println(errMsg)
			errMsgs += errMsg
		}
		return &info.ActionList{}, nil, errors.New(errMsgs)
	}

	return simcfg, gcsl, nil
}

func setupAggregators(simcfg *info.ActionList) ([]agg.Aggregator, error) {
	var aggregators []agg.Aggregator
	enabled := simcfg.Settings.CollectStats
	for _, aggregator := range agg.Aggregators() {
		if len(enabled) > 0 && !slices.Contains(enabled, aggregator.Name) {
			continue
		}
		a, err := aggregator.New(simcfg)
		if err != nil {
			return nil, err
		}
		aggregators = append(aggregators, a)
	}
	return aggregators, nil
}

func cryptoRandSeed() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}
