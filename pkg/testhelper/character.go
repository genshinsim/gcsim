package testhelper

import (
	_ "embed"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/prototext"
)

//go:embed test_char_data.pb
var pbData []byte
var base *model.AvatarData

func init() {
	base = &model.AvatarData{}
	err := prototext.Unmarshal(pbData, base)
	if err != nil {
		panic(err)
	}
}

type Character struct {
	*character.CharWrapper
}

func (c *Character) Snapshot(a *combat.AttackInfo) combat.Snapshot { return combat.Snapshot{} }
func (c *Character) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	return true, action.NoFailure
}
func (c *Character) NextQueueItemIsValid(a action.Action, p map[string]int) error {
	return nil
}
func (c *Character) ActionStam(a action.Action, p map[string]int) float64 { return 0 }
func (c *Character) ReduceActionCooldown(a action.Action, v int)          {}
func (c *Character) ResetActionCooldown(a action.Action)                  {}
func (c *Character) Cooldown(a action.Action) int                         { return 0 }
func (c *Character) SetCDWithDelay(a action.Action, dur, delay int)       {}
func (c *Character) Charges(a action.Action) int                          { return 1 }
func (c *Character) SetCD(a action.Action, dur int)                       {}
func (c *Character) Init() error                                          { return nil }
func (c *Character) Data() *model.AvatarData                              { return base }
func (c *Character) CurrentHPRatio() float64                              { return 0 }
func (c *Character) CurrentHP() float64                                   { return 0 }
func (c *Character) CurrentHPDebt() float64                               { return 0 }
func (c *Character) SetHPByAmount(float64)                                {}
func (c *Character) SetHPByRatio(float64)                                 {}
func (c *Character) ModifyHPByAmount(float64)                             {}
func (c *Character) ModifyHPByRatio(float64)                              {}
func (c *Character) ModifyHPDebtByAmount(float64)                         {}
func (c *Character) ModifyHPDebtByRatio(float64)                          {}
func (c *Character) Heal(*info.HealInfo) (float64, float64)               { return 0, 0 }
func (c *Character) Drain(*info.DrainInfo) float64                        { return 0 }

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := Character{}
	c.CharWrapper = w
	w.Character = &c
	return nil
}
