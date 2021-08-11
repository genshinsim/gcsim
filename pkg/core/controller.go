package core

import (
	"go.uber.org/zap"
	"math/rand"
)

type ShieldHandler interface {
	AddShield(shd Shield)
	IsShielded() bool
	GetShield(t ShieldType) Shield
	AddShieldBonus(f func() float64)
}

type ConstructHandler interface {
	NewConstruct(c Construct, refresh bool)
	NewNoLimitCons(c Construct, refresh bool)
	ConstructCount() int
	ConstructCountType(t GeoConstructType) int
	Destroy(key int) bool
	HasConstruct(key int) bool
}

type HPHandler interface {

}

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
	ConstructCtrl ConstructHandler

	//shields
	ShieldCtrl ShieldHandler
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
