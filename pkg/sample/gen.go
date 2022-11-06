package sample

import (
	"encoding/json"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/simulation"
	"github.com/genshinsim/gcsim/pkg/simulator"
)

// GenerateSampleWithSeed will run one simulation with debug enabled using the given seed and output
// the debug log. Used for generating debug for min/max runs
func GenerateSampleWithSeed(cfg string, seed uint64) (Sample, error) {
	simcfg, err := simulator.Parse(cfg)
	if err != nil {
		return Sample{}, err
	}

	c, err := simulation.NewCore(int64(seed), true, simcfg)
	if err != nil {
		return Sample{}, err
	}

	//create a new simulation and run
	s, err := simulation.New(simcfg, c)
	if err != nil {
		return Sample{}, err
	}
	_, err = s.Run()
	if err != nil {
		return Sample{}, err
	}

	//capture the log
	logs, err := c.Log.Dump()
	if err != nil {
		return Sample{}, err
	}

	// TODO: Log.Dump() should not marshal the data. Embedding json as a string in json is just bad
	var events []map[string]interface{}
	if err := json.Unmarshal(logs, &events); err != nil {
		return Sample{}, err
	}

	chars, err := simulator.GenerateCharacterDetails(simcfg)
	if err != nil {
		return Sample{}, err
	}

	sample := Sample{
		Config:           cfg,
		CharacterDetails: chars,
		TargetDetails:    simcfg.Targets,
		Seed:             strconv.FormatUint(seed, 10),
		Logs:             events,
	}

	sample.Logs = events
	return sample, err
}
