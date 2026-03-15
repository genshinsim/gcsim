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
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterCharFunc(keys.Nefer, NewChar)
}

const (
	shadowDanceKey          = "nefer-shadow-dance"
	slitherKey              = "nefer-slither"
	phantasmChargeHoldKey   = "nefer-phantasm-charge-hold"
	seedWindowKey           = "nefer-seed-window"
	veilEMBuffKey           = "nefer-veil-em-buff"
	seedAbsorbRadius        = 6
	phantasmChargesPerSkill = 3
)

type char struct {
	*tmpl.Character
	ascendantGleam    bool
	veilstacks        int
	maxVeilStacks     int
	slitherSrc        int
	phantasmCharges   int
	phantasmChargeSrc int
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
	c.swapResetInit()
	c.lunarbloomInit()
	c.p1Init()
	c.c6Init()
	return nil
}

func (c *char) slitherActive() bool {
	return c.StatusIsActive(slitherKey)
}

func (c *char) phantasmChargeHoldActive() bool {
	return c.StatusIsActive(phantasmChargeHoldKey)
}

func (c *char) hasPhantasmCharge() bool {
	return c.phantasmCharges > 0
}

func (c *char) canTriggerPhantasm() bool {
	return c.StatusIsActive(shadowDanceKey) && c.hasPhantasmCharge() && c.Core.Player.VerdantDew() > 0
}

func (c *char) phantasmActive() bool {
	return c.phantasmEndFrame > c.Core.F
}

func (c *char) chargeLoopActive(src int) bool {
	return c.phantasmChargeSrc == src && c.phantasmChargeHoldActive()
}

func (c *char) chargeLoopInterrupted(src int) bool {
	if !c.chargeLoopActive(src) {
		return true
	}
	if c.StatusIsActive(shadowDanceKey) && c.Core.Player.Active() == c.Index() && c.Core.Player.CurrentState() == action.ChargeAttackState {
		return false
	}
	c.clearPhantasmChargeLoop()
	return true
}

func (c *char) slitherLoopInterrupted(src int) bool {
	if c.chargeLoopInterrupted(src) {
		return true
	}
	return c.slitherSrc != src || !c.slitherActive()
}

func (c *char) clearSlither() {
	c.DeleteStatus(slitherKey)
	c.slitherSrc = 0
}

func (c *char) clearPhantasmChargeLoop() {
	c.clearSlither()
	c.DeleteStatus(phantasmChargeHoldKey)
	c.phantasmChargeSrc = 0
	c.phantasmEndFrame = 0
}

func (c *char) swapResetInit() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		if prev != c.Index() {
			return
		}
		c.clearPhantasmChargeLoop()
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
	prev := c.veilstacks
	c.veilstacks = min(c.veilstacks+count, c.maxVeilStacks)
	c.applyVeilThresholdBuff(prev, c.veilstacks)
}

func (c *char) consumeVeilStacks() int {
	stacks := c.veilstacks
	c.veilstacks = 0
	return stacks
}

func (c *char) applyVeilThresholdBuff(prev, next int) {
	if next <= prev {
		return
	}

	amount := 0.0
	if c.Base.Cons >= 2 && next >= 5 && prev < 5 {
		amount = 200
	} else if next >= 3 && prev < 3 {
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
