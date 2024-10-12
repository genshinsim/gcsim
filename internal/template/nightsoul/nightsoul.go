package nightsoul

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const NightsoulBlessingStatus = "nightsoul-blessing"

type State struct {
	char            *character.CharWrapper
	c               *core.Core
	nightsoulPoints float64

	MaxPoints float64
}

func New(c *core.Core, char *character.CharWrapper) *State {
	t := &State{
		char:      char,
		c:         c,
		MaxPoints: -1.0, // no limits
	}
	return t
}

func (s *State) EnterBlessing(amount float64) {
	s.nightsoulPoints = amount
	s.char.AddStatus(NightsoulBlessingStatus, -1, true)
	s.c.Log.NewEvent("enter nightsoul blessing", glog.LogCharacterEvent, s.char.Index).
		Write("points", s.nightsoulPoints)
}

func (s *State) ExitBlessing() {
	s.char.DeleteStatus(NightsoulBlessingStatus)
	s.c.Log.NewEvent("exit nightsoul blessing", glog.LogCharacterEvent, s.char.Index)
}

func (s *State) HasBlessing() bool {
	return s.char.StatusIsActive(NightsoulBlessingStatus)
}

func (s *State) GeneratePoints(amount float64) {
	prevPoints := s.nightsoulPoints
	s.nightsoulPoints += amount
	s.clampPoints()
	s.c.Events.Emit(event.OnNightsoulGenerate, s.char.Index, amount)
	s.c.Log.NewEvent("generate nightsoul points", glog.LogCharacterEvent, s.char.Index).
		Write("previous points", prevPoints).
		Write("amount", amount).
		Write("final", s.nightsoulPoints)
}

func (s *State) ConsumePoints(amount float64) {
	prevPoints := s.nightsoulPoints
	s.nightsoulPoints -= amount
	s.clampPoints()
	s.c.Events.Emit(event.OnNightsoulConsume, s.char.Index, amount)
	s.c.Log.NewEvent("consume nightsoul points", glog.LogCharacterEvent, s.char.Index).
		Write("previous points", prevPoints).
		Write("amount", amount).
		Write("final", s.nightsoulPoints)
}

func (s *State) clampPoints() {
	if s.MaxPoints > 0 && s.nightsoulPoints > s.MaxPoints {
		s.nightsoulPoints = s.MaxPoints
	} else if s.nightsoulPoints < 0 {
		s.nightsoulPoints = 0
	}
}

func (s *State) Points() float64 {
	return s.nightsoulPoints
}
