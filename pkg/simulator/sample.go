package simulator

import (
	"encoding/json"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulation"
	"google.golang.org/protobuf/types/known/structpb"
)

// GenerateSampleWithSeed will run one simulation with debug enabled using the given seed and output
// the debug log. Used for generating debug for min/max runs
func GenerateSampleWithSeed(cfg string, seed uint64, opts Options) (*model.Sample, error) {
	simcfg, err := Parse(cfg)
	if err != nil {
		return &model.Sample{}, err
	}

	c, err := simulation.NewCore(int64(seed), true, simcfg)
	if err != nil {
		return &model.Sample{}, err
	}

	//create a new simulation and run
	s, err := simulation.New(simcfg, c)
	if err != nil {
		return &model.Sample{}, err
	}
	_, err = s.Run()
	if err != nil {
		return &model.Sample{}, err
	}

	//capture the log
	logs, err := c.Log.Dump()
	if err != nil {
		return &model.Sample{}, err
	}

	// TODO: Log.Dump() should not marshal the data. Embedding json as a string in json is just bad
	var events []map[string]interface{}
	if err := json.Unmarshal(logs, &events); err != nil {
		return &model.Sample{}, err
	}

	chars, err := GenerateCharacterDetails(simcfg)
	if err != nil {
		return &model.Sample{}, err
	}

	sample := &model.Sample{
		Config:           cfg,
		InitialCharacter: simcfg.InitialChar.String(),
		CharacterDetails: chars,
		Seed:             strconv.FormatUint(seed, 10),
		TargetDetails:    make([]*model.Enemy, len(simcfg.Targets)),
		Logs:             make([]*structpb.Struct, len(events)),
	}

	sample.TargetDetails = make([]*model.Enemy, len(simcfg.Targets))
	for i, target := range simcfg.Targets {
		resist := make(map[string]float64)
		for k, v := range target.Resist {
			resist[k.String()] = v
		}

		sample.TargetDetails[i] = &model.Enemy{
			Level:  int32(target.Level),
			HP:     target.HP,
			Resist: resist,
			Pos: &model.Coord{
				X: target.Pos.X,
				Y: target.Pos.Y,
				R: target.Pos.R,
			},
			ParticleDropThreshold: target.ParticleDropThreshold,
			ParticleDropCount:     target.ParticleDropCount,
			ParticleElement:       target.ParticleElement.String(),
		}
	}

	for i, event := range events {
		es, err := structpb.NewStruct(event)
		if err != nil {
			return &model.Sample{}, err
		}
		sample.Logs[i] = es
	}
	return sample, err
}
