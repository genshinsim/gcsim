package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type Character struct {
	*character.CharWrapper
	Core                   *core.Core
	ActionCD               []int
	cdQueueWorkerStartedAt []int
	cdCurrentQueueWorker   []*func()
	cdQueue                [][]int
	AvailableCDCharge      []int
	additionalCDCharge     []int

	// hp
	currentHPRatio float64
	currentHPDebt  float64
}

func NewWithWrapper(c *core.Core, w *character.CharWrapper) *Character {
	r := New(c)
	r.CharWrapper = w
	return r
}

func New(c *core.Core) *Character {
	t := &Character{
		Core:                   c,
		ActionCD:               make([]int, action.EndActionType),
		cdQueueWorkerStartedAt: make([]int, action.EndActionType),
		cdCurrentQueueWorker:   make([]*func(), action.EndActionType),
		cdQueue:                make([][]int, action.EndActionType),
		AvailableCDCharge:      make([]int, action.EndActionType),
		additionalCDCharge:     make([]int, action.EndActionType),
	}

	for i := 0; i < len(t.cdQueue); i++ {
		t.cdQueue[i] = make([]int, 0, 4)
		t.AvailableCDCharge[i] = 1
	}

	return t
}

func (c *Character) Snapshot(a *combat.AttackInfo) combat.Snapshot {
	s := combat.Snapshot{
		CharLvl:     c.Base.Level,
		ActorEle:    c.Base.Element,
		BaseAtk:     c.Base.Atk + c.Weapon.BaseAtk,
		BaseDef:     c.Base.Def,
		BaseHP:      c.Base.HP,
		SourceFrame: c.Core.F,
	}

	var evt glog.Event
	var debug []interface{}

	if c.Core.Flags.LogDebug {
		evt = c.Core.Log.NewEvent(a.Abil, glog.LogSnapshotEvent, c.Index).
			Write("abil", a.Abil).
			Write("mult", a.Mult).
			Write("ele", a.Element.String()).
			Write("durability", float64(a.Durability)).
			Write("icd_tag", a.ICDTag).
			Write("icd_group", a.ICDGroup)
	}

	// snapshot the stats
	s.Stats, debug = c.Stats()

	// check infusion
	var inf attributes.Element
	if !a.IgnoreInfusion {
		inf = c.Core.Player.Infused(c.Index, a.AttackTag)
		if inf != attributes.NoElement {
			a.Element = inf
		}
	}

	// check if we need to log
	if c.Core.Flags.LogDebug {
		evt.WriteBuildMsg(debug...)
		evt.Write("final_stats", attributes.PrettyPrintStatsSlice(s.Stats[:]))
		if inf != attributes.NoElement {
			evt.Write("infused_ele", inf.String())
		}
		s.Logs = debug
	}
	return s
}

func (c *Character) ResetNormalCounter() {
	c.NormalCounter = 0
}

func (c *Character) AdvanceNormalIndex() {
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

func (c *Character) NextNormalCounter() int {
	return c.NormalCounter + 1
}

func (c *Character) Attack(map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action attack not implemented", c.CharWrapper.Base.Key)
}
func (c *Character) Aimed(map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action aimed not implemented", c.CharWrapper.Base.Key)
}
func (c *Character) ChargeAttack(map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action charge not implemented", c.CharWrapper.Base.Key)
}
func (c *Character) HighPlungeAttack(map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action high_plunge not implemented", c.CharWrapper.Base.Key)
}
func (c *Character) LowPlungeAttack(map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action low_plunge not implemented", c.CharWrapper.Base.Key)
}
func (c *Character) Skill(map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action skill not implemented", c.CharWrapper.Base.Key)
}
func (c *Character) Burst(map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action burst not implemented", c.CharWrapper.Base.Key)
}

func (c *Character) NextQueueItemIsValid(_ keys.Char, a action.Action, p map[string]int) error {
	if a == action.ActionCharge {
		switch c.Weapon.Class {
		case info.WeaponClassSword, info.WeaponClassSpear:
			// cannot do charge on most sword/polearm characters without attack beforehand
			if c.Core.Player.LastAction.Type != action.ActionAttack {
				return player.ErrInvalidChargeAction
			}
		}
	}
	return nil
}
