package core

import (
	"go.uber.org/zap"
	"math/rand"
)

type Controller struct {
	//control
	F int //frame
	Rand *rand.Rand
	Log *zap.SugaredLogger

	//player
	Stam int

	//characters
	ActiveChar int
	Chars []Character

	//enemy
	Targets []Target

	//constructs
	Constructs []Construct
	ConstructsNoLimit []Construct

	//shields
	Shields []Shield
}

func New(cfg ...func(*Controller) error) (*Controller, error) {
	c := &Controller{}
	for _, f := range cfg {
		err := f(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}
