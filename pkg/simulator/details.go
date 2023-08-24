package simulator

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

func GenerateCharacterDetails(cfg *info.ActionList) ([]*model.Character, error) {
	cpy := cfg.Copy()

	c, err := simulation.NewCore(CryptoRandSeed(), false, cpy)
	if err != nil {
		return nil, err
	}
	//create a new simulation and run
	//TODO: nil shoudl be fine here
	sim, err := simulation.New(cpy, nil, c)
	if err != nil {
		return nil, err
	}

	return sim.CharacterDetails(), nil
}
