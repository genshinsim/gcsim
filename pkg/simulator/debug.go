package simulator

import (
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

//GenerateDebugLogWithSeed will run one simulation with debug enabled using the given seed and output
//the debug log. Used for generating debug for min/max runs
func GenerateDebugLogWithSeed(cfg *ast.ActionList, seed int64) (string, error) {
	cpy := cfg.Copy()

	c, err := simulation.NewCore(seed, true, cpy)
	if err != nil {
		return "", err
	}
	//create a new simulation and run
	s, err := simulation.New(cpy, c)
	if err != nil {
		return "", err
	}
	_, err = s.Run()
	if err != nil {
		return "", err
	}
	//capture the log
	out, err := c.Log.Dump()
	return string(out), err
}

//GenerateDebugLog will run one simulation with debug enabled using a random seed
func GenerateDebugLog(cfg *ast.ActionList) (string, error) {
	return GenerateDebugLogWithSeed(cfg, cryptoRandSeed())
}
