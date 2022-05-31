package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type EnergyEvent struct {
	Active    bool
	Once      bool //how often
	Start     int
	End       int
	Particles int
}

type HurtEvent struct {
	Active bool
	Once   bool //how often
	Start  int  //
	End    int
	Min    float64
	Max    float64
	Ele    attributes.Element
}
