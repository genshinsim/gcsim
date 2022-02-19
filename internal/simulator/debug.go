package simulator

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

//GenerateDebugLogWithSeed will run one simulation with debug enabled using the given seed and output
//the debug log. Used for generating debug for min/max runs
func GenerateDebugLogWithSeed(cfg core.SimulationConfig, seed int64) (string, error) {

	c := simulation.NewCore(seed, true, cfg.Settings)
	//create a new simulation and run
	s, err := simulation.New(cfg, c)
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
func GenerateDebugLog(cfg core.SimulationConfig) (string, error) {
	return GenerateDebugLogWithSeed(cfg, cryptoRandSeed())
}
