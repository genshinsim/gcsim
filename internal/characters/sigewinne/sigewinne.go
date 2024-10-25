package sigewinne

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/sourcewaterdroplet"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Sigewinne, NewChar)
}

type char struct {
	*tmpl.Character

	skillAttackInfo combat.AttackInfo
	skillSnapshot   combat.Snapshot

	particleGenerated   bool
	lastSummonSrc       int
	bubbleHitLimit      int
	currentBubbleTier   int
	collectedHpDebt     float32
	burstEarlyCancelled bool
	tickAnimLength      int
	burstMaxDuration    int
	burstStartF         int
	lastSwap            int
	chargeAi            combat.AttackInfo
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3
	c.HasArkhe = true

	c.bubbleHitLimit = 5
	c.currentBubbleTier = 0
	c.collectedHpDebt = 0
	c.burstEarlyCancelled = false
	c.burstMaxDuration = 241 - chargeBurstDur
	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Ascension >= 1 {
		c.a1()
	}
	if c.Base.Ascension >= 4 {
		c.a4()
	}

	if c.Base.Cons >= 1 {
		c.bubbleHitLimit += 3
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 4 {
		c.burstMaxDuration = 425 - chargeBurstDur
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.bubbleTierDamageMod()
	c.energyBondClearMod()
	c.onSwap()
	return nil
}

func (c *char) getSourcewaterDroplets() []*sourcewaterdroplet.Gadget {
	player := c.Core.Combat.Player()

	// Used Neuvillette's droplet tracking
	// TODO: check if true for Sigewinne
	segment := combat.NewCircleHitOnTargetFanAngle(player, nil, 14, 80)
	rect := combat.NewBoxHitOnTarget(player, geometry.Point{Y: -7}, 8, 18)

	droplets := make([]*sourcewaterdroplet.Gadget, 0)
	for _, g := range c.Core.Combat.Gadgets() {
		droplet, ok := g.(*sourcewaterdroplet.Gadget)
		if !ok {
			continue
		}
		if !droplet.IsWithinArea(rect) && !droplet.IsWithinArea(segment) {
			continue
		}
		droplets = append(droplets, droplet)
	}

	return droplets
}

// used for early Burst cancel swap cd calculation
func (c *char) onSwap() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		next := args[1].(int)
		if next != c.Index {
			return false
		}
		c.lastSwap = c.Core.F
		return false
	}, "sigewinne-swap")
}

func (c *char) consumeDroplet(g *sourcewaterdroplet.Gadget) {
	g.Kill()
	c.ModifyHPDebtByAmount(c.MaxHP() * BoLPctPerDroplet)
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 13
	}
	if k == model.AnimationYelanN0StartDelay {
		return 5
	}
	return c.Character.AnimationStartDelay(k)
}
