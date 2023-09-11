package testhelper

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type Character struct {
	*character.CharWrapper
}

func (c *Character) Snapshot(a *combat.AttackInfo) combat.Snapshot { return combat.Snapshot{} }
func (c *Character) ActionReady(a action.Action, p map[string]int) (bool, action.ActionFailure) {
	return true, action.NoFailure
}
func (c *Character) ActionStam(a action.Action, p map[string]int) float64 { return 0 }
func (c *Character) ReduceActionCooldown(a action.Action, v int)          {}
func (c *Character) ResetActionCooldown(a action.Action)                  {}
func (c *Character) Cooldown(a action.Action) int                         { return 0 }
func (c *Character) SetCDWithDelay(a action.Action, dur int, delay int)   {}
func (c *Character) Charges(a action.Action) int                          { return 1 }
func (c *Character) SetCD(a action.Action, dur int)                       {}
func (c *Character) Init() error                                          { return nil }

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := Character{}
	c.CharWrapper = w
	w.Character = &c
	return nil
}
