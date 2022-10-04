package simulator

import (
	"encoding/json"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

// GenerateDebugLogWithSeed will run one simulation with debug enabled using the given seed and output
// the debug log. Used for generating debug for min/max runs
func GenerateDebugLogWithSeed(cfg *ast.ActionList, seed int64) ([]map[string]interface{}, error) {
	cpy := cfg.Copy()

	c, err := simulation.NewCore(seed, true, cpy)
	if err != nil {
		return nil, err
	}
	//create a new simulation and run
	s, err := simulation.New(cpy, c)
	if err != nil {
		return nil, err
	}
	_, err = s.Run()
	if err != nil {
		return nil, err
	}
	//capture the log
	out, err := c.Log.Dump()
	if err != nil {
		return nil, err
	}

	// TODO: Log.Dump() should not marshal the data. Embedding json as a string in json is just bad
	var events []map[string]interface{}
	if err := json.Unmarshal(out, &events); err != nil {
		return nil, err
	}

	return events, err
}

// GenerateDebugLog will run one simulation with debug enabled using a random seed
func GenerateDebugLog(cfg *ast.ActionList) ([]map[string]interface{}, error) {
	return GenerateDebugLogWithSeed(cfg, CryptoRandSeed())
}
