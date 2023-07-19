package simulator

import (
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

func GenerateCharacterDetails(cfg *ast.ActionList) ([]*model.Character, error) {
	cpy := cfg.Copy()

	c, err := simulation.NewCore(CryptoRandSeed(), false, cpy)
	if err != nil {
		return nil, err
	}
	//create a new simulation and run
	sim, err := simulation.New(cpy, c)
	if err != nil {
		return nil, err
	}

	return sim.CharacterDetails(), nil
}
