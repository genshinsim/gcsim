package nefer

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/stacks"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterCharFunc(keys.Nefer, NewChar)
}

const (
	shadowDanceKey          = "nefer-shadow-dance"
	slitherKey              = "nefer-slither"
	seedWindowKey           = "nefer-seed-window"
	veilEMBuffKey           = "nefer-veil-em-buff"
	seedAbsorbRadius        = 6
	veilBaseDuration        = 9 * 60
	c2VeilDurationBonus     = 5 * 60
	phantasmChargesPerSkill = 3
)

type char struct {
	*tmpl.Character
	ascendantGleam    bool
	veilstacks        int
	veilTracker       *stacks.MultipleRefreshNoRemove
	maxVeilStacks     int
	phantasmCharges   int
	chargeRoute       chargeRouteState
}

type chargeRouteState struct {
	src               int
	slitherSrc        int
	releaseStartFrame int
	phantasmStartFrame int
	phantasmEndFrame  int
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
	c.maxVeilStacks = 3

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.ascendantGleam = c.Core.Player.GetMoonsignLevel() >= 2
	if c.Base.Cons >= 2 {
		c.maxVeilStacks = 5
	}
	c.initVeilTracker()
	c.swapResetInit()
	c.lunarbloomInit()
	c.p1Init()
	c.p2Init()
	c.c4Init()
	c.c6Init()
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
	prev := c.currentVeilStacks()
	for range count {
		c.veilTracker.Add(c.veilStackDuration())
		next := c.currentVeilStacks()
		c.applyVeilThresholdBuff(prev, next, 1)
		prev = next
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
		return 0
	}
	return float64(stacks) * 0.08
}

func (c *char) initVeilTracker() {
	c.veilTracker = stacks.NewMultipleRefreshNoRemove(c.maxVeilStacks, c.QueueCharTask, &c.Core.F)
	c.veilstacks = 0
}

func (c *char) currentVeilStacks() int {
	if c.veilTracker == nil {
		return c.veilstacks
	}
	c.veilstacks = c.veilTracker.Count()
	return c.veilstacks
}

func (c *char) veilStackDuration() int {
	dur := veilBaseDuration
	if c.Base.Cons >= 2 {
		dur += c2VeilDurationBonus
	}
	return dur
}

func (c *char) applyVeilThresholdBuff(prev, next, count int) {
	if count <= 0 {
		return
	}

	triggerFive := c.Base.Cons >= 2 && next >= 5 && (prev < 5 || prev == 5)
	triggerThree := next >= 3 && (prev < 3 || prev == 3)

	amount := 0.0
	if triggerFive {
		amount = 200
	} else if triggerThree {
		amount = 100
	}
	if amount <= 0 {
		return
	}

	buff := make([]float64, attributes.EndStatType)
	buff[attributes.EM] = amount
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(veilEMBuffKey, 8*60),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			return buff
		},
	})
}
