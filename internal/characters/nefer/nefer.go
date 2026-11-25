package nefer

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/stacks"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	shadowDanceKey   = "nefer-shadow-dance"
	slitherKey       = "nefer-slither"
	seedWindowKey    = "nefer-seed-window"
	veilEMBuffKey    = "nefer-veil-em-buff"
	seedAbsorbRadius = 6
)

type char struct {
	*tmpl.Character
	veilTracker     *stacks.MultipleRefreshNoRemove
	maxVeilStacks   int
	phantasmCharges int
	a4BuffSrc       int
	a4BuffAmt       float64
	chargeRoute     chargeRouteState
}

type chargeRouteState struct {
	src                int
	slitherSrc         int
	releaseStartFrame  int
	phantasmStartFrame int
	phantasmEndFrame   int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5
	c.Moonsign = 1
	c.SetNumCharges(action.ActionSkill, 2)

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.swapResetInit()
	c.lunarbloomInit()
	c.a1Init()
	c.a4Init()
	c.c2VeilStackInit()
	c.c4Init()
	c.c6Init()
	c.initVeilTracker()
	return nil
}

func (c *char) slitherActive() bool {
	return c.StatusIsActive(slitherKey)
}

func (c *char) hasPhantasmCharge() bool {
	return c.phantasmCharges > 0
}

func (c *char) canTriggerPhantasm() bool {
	return c.StatusIsActive(shadowDanceKey) && c.hasPhantasmCharge() && c.Core.Player.VerdantDew() > 0
}

func (c *char) phantasmActive() bool {
	return c.chargeRoute.phantasmEndFrame > c.Core.F
}

func (c *char) chargeRouteActive(src int) bool {
	return c.chargeRoute.src == src && c.Core.Player.Active() == c.Index() && c.Core.Player.CurrentState() == action.ChargeAttackState
}

func (c *char) chargeRouteInterrupted(src int) bool {
	if c.chargeRouteActive(src) {
		return false
	}
	c.clearPhantasmChargeLoop()
	return true
}

func (c *char) slitherLoopInterrupted(src int) bool {
	if c.chargeRouteInterrupted(src) {
		return true
	}
	return c.chargeRoute.slitherSrc != src || !c.slitherActive()
}

func (c *char) clearSlither() {
	c.DeleteStatus(slitherKey)
	c.chargeRoute.slitherSrc = 0
}

func (c *char) clearPhantasmChargeLoop() {
	c.clearSlither()
	c.chargeRoute = chargeRouteState{}
}

func (c *char) swapResetInit() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		if prev != c.Index() {
			return
		}
		c.clearPhantasmChargeLoop()
		c.DeleteStatus(shadowDanceKey)
		c.phantasmCharges = 0
	}, "nefer-swap-reset")
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationYelanN0StartDelay:
		return 10
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) addVeilStacks(count int) {
	if count <= 0 {
		return
	}
	for range count {
		c.veilTracker.Add(9*60 + c.c2VeilDurationBonus())
		c.applyVeilThresholdBuff(c.currentVeilStacks())
	}
}

func (c *char) consumeVeilStacks() int {
	count := c.currentVeilStacks()
	c.initVeilTracker()
	return count
}

func (c *char) phantasmVeilBonus() float64 {
	stacks := c.currentVeilStacks()
	if stacks <= 0 {
		return 1
	}
	return 1 + float64(stacks)*0.08
}

func (c *char) initVeilTracker() {
	c.veilTracker = stacks.NewMultipleRefreshNoRemove(c.maxVeilStacks, c.QueueCharTask, &c.Core.F)
}

func (c *char) currentVeilStacks() int {
	if c.veilTracker == nil {
		return 0
	}

	return c.veilTracker.Count()
}

func (c *char) veilStackFull(count int) bool {
	return count >= c.maxVeilStacks
}

func (c *char) applyVeilThresholdBuff(count int) {
	buff := make([]float64, attributes.EndStatType)
	key := veilEMBuffKey
	switch {
	case c.c2VeilStackFull(count):
		buff[attributes.EM] = 200
		key = veilEMBuffKey + "-c2"
	case c.veilStackFull(count):
		buff[attributes.EM] = 100
	default:
		return
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(key, 8*60),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			return buff
		},
	})
}
