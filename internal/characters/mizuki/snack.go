package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const (
	snackDmgName              = "Munen Shockwave"
	snackHealName             = "Snack Pick-Up"
	snackDurability           = 25
	snackDmgRadius            = 4
	snackHealTriggerHpRatio   = 0.7
	snackDuration             = 4 * 60
	snackSize                 = 2.5
	snackSizeMizukiMultiplier = 2.1 // Assumption
	snackCantTriggerDuration  = 0.3 * 60
)

type snack struct {
	*gadget.Gadget
	char         *char
	attackInfo   info.AttackInfo
	snapshot     info.Snapshot
	pattern      info.AttackPattern
	allowPickupF int
}

func newSnack(c *char, pos info.Point) *snack {
	p := &snack{
		char: c,
		attackInfo: info.AttackInfo{
			ActorIndex:   c.Index(),
			Abil:         snackDmgName,
			AttackTag:    attacks.AttackTagElementalBurst,
			ICDTag:       attacks.ICDTagElementalBurst,
			ICDGroup:     attacks.ICDGroupDefault,
			StrikeType:   attacks.StrikeTypeDefault,
			Element:      attributes.Anemo,
			Durability:   snackDurability,
			Mult:         snackDMG[c.TalentLvlBurst()],
			HitlagFactor: 0.05,
		},
		pattern:      combat.NewCircleHitOnTarget(pos, nil, snackDmgRadius),
		allowPickupF: c.Core.F + snackCantTriggerDuration,
	}
	p.snapshot = c.Snapshot(&p.attackInfo)

	// we increase snack size to make sure we get it when mizuki is in dreamdrifter state
	// because mizuki's pickup range is increased while in this state.
	// https://docs.google.com/spreadsheets/d/1UU0EVPBatEndl4GRZyIs8Ix8O3kcZUDAwOHqM8_jQJw/edit?gid=339012102#gid=339012102
	p.Gadget = gadget.New(c.Core, pos, snackSize*snackSizeMizukiMultiplier, info.GadgetTypYumemiSnack)
	p.Duration = snackDuration
	c.Core.Combat.AddGadget(p)

	p.CollidableTypes[info.TargettablePlayer] = true
	p.OnExpiry = func() {
		p.explode()
		p.Core.Log.NewEvent("Snack exploded by itself", glog.LogCharacterEvent, c.Index())
	}
	p.OnCollision = func(target info.Target) {
		if _, ok := target.(*avatar.Player); !ok {
			return
		}
		if p.Core.F < p.allowPickupF {
			return
		}

		// default size is increased. The increased size is only valid for mizuki in dreamdrifter state,
		// so check for collision with actual size if this is not the case
		if !c.StatusIsActive(dreamDrifterStateKey) && !p.collidesWithActiveCharacterDefaultSize() {
			return
		}
		p.onPickedUp()
	}

	p.Core.Log.NewEvent("Snack spawned", glog.LogCharacterEvent, c.Index()).
		Write("x", pos.X).
		Write("y", pos.Y)
	return p
}

func (p *snack) collidesWithActiveCharacterDefaultSize() bool {
	defaultSize := combat.NewCircleHitOnTarget(p.Gadget, nil, snackSize)
	return p.Core.Combat.Player().IsWithinArea(defaultSize)
}

func (p *snack) onPickedUp() {
	var heal bool
	var dmg bool

	mizuki := p.char
	activeChar := p.Core.Player.ActiveChar()

	// C4 triggers both DMG and Heal
	if mizuki.Base.Cons >= 4 {
		dmg = true
		heal = true
	} else {
		// Heals active char if is bellow 70% hp otherwise deals DMG
		dmg = activeChar.CurrentHP() > (activeChar.MaxHP() * snackHealTriggerHpRatio)
		heal = !dmg
	}

	p.Core.Log.NewEvent("Picked up snack", glog.LogCharacterEvent, activeChar.Index()).
		Write("heal", heal).
		Write("dmg", dmg)

	if dmg {
		p.explode()
	}

	if heal {
		// Heals double the amount on Mizuki
		healMultiplier := 1.0
		if activeChar.Index() == mizuki.Index() {
			healMultiplier = 2.0
		}
		mizuki.Core.Player.Heal(info.HealInfo{
			Caller:  mizuki.Index(),
			Target:  activeChar.Index(),
			Message: snackHealName,
			Src:     ((mizuki.Stat(attributes.EM) * snackHealEM[mizuki.TalentLvlBurst()]) + snackHealFlat[mizuki.TalentLvlBurst()]) * healMultiplier,
			Bonus:   mizuki.Stat(attributes.Heal),
		})
	}

	// C4 restores 5 energy to mizuki up to 4 times
	mizuki.c4()

	p.Kill()
}

func (p *snack) explode() {
	p.Core.QueueAttackWithSnap(p.attackInfo, p.snapshot, p.pattern, 0)
}

func (p *snack) HandleAttack(atk *info.AttackEvent) float64 {
	// only collisions with the player can affect this or if it expires
	return 0
}

func (p *snack) SetDirection(trg info.Point) {}
func (p *snack) SetDirectionToClosestEnemy() {}
func (p *snack) CalcTempDirection(trg info.Point) info.Point {
	return info.DefaultDirection()
}
