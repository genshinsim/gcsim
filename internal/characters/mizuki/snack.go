package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const (
	snackDuration             = 4 * 60
	snackSize                 = 2
	snackSizeMizukiMultiplier = 1.75 // Assumption
)

type snack struct {
	*gadget.Gadget
	char       *char
	attackInfo combat.AttackInfo
	snapshot   combat.Snapshot
	pattern    combat.AttackPattern
}

func newSnack(c *char, pos geometry.Point) *snack {
	p := &snack{
		char: c,
		attackInfo: combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       snackDmgName,
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: snackDurability,
			PoiseDMG:   snackPoise,
			Mult:       snackDMG[c.TalentLvlBurst()],
		},
		pattern: combat.NewCircleHitOnTarget(pos, nil, snackRadius),
	}
	p.snapshot = c.Snapshot(&p.attackInfo)

	// we increase snack size to make sure we get it when mizuki is in dreamdrifter state
	p.Gadget = gadget.New(c.Core, pos, snackSize*snackSizeMizukiMultiplier, combat.GadgetTypYumemiSnack)
	p.Gadget.Duration = snackDuration
	c.Core.Combat.AddGadget(p)

	p.Gadget.CollidableTypes[targets.TargettablePlayer] = true
	p.Gadget.OnExpiry = func() {
		p.explodeWithHitmark(snackHitmark)
		p.char.Core.Log.NewEvent("Snack exploded by itself", glog.LogCharacterEvent, c.Index)
	}
	p.Gadget.OnCollision = func(target combat.Target) {
		if _, ok := target.(*avatar.Player); !ok {
			return
		}

		// default size is increased. The increased size is only valid for mizuki in dreamdrifter state,
		// so check for collision with actual size if this is not the case
		collidesDefault := p.collidesWithActiveCharacterDefaultSize()
		if c.Core.Player.ActiveChar().Index != c.Index {
			if !collidesDefault {
				return
			}
		} else {
			if !c.StatusIsActive(dreamDrifterStateKey) && !collidesDefault {
				return
			}
		}
		p.onPickedUp()
	}

	p.char.Core.Log.NewEvent("Snack spawned", glog.LogCharacterEvent, c.Index).
		Write("x", pos.X).
		Write("y", pos.Y)
	return p
}

func (p *snack) collidesWithActiveCharacterDefaultSize() bool {
	return p.char.Core.Combat.Player().WillCollide(geometry.NewCircle(p.Gadget.Pos(), snackSize, geometry.DefaultDirection(), 360))
}

func (p *snack) onPickedUp() {
	heal := false
	dmg := false

	mizuki := p.char
	activeChar := p.char.Core.Player.ActiveChar()

	// C4 triggers both DMG and Heal
	if mizuki.Base.Cons >= 4 {
		dmg = true
		heal = true
	} else {
		// Heals active char if is bellow 70% hp otherwise deals DMG
		dmg = activeChar.CurrentHP() > (activeChar.MaxHP() * snackHealTriggerHpRatio)
		heal = !dmg
	}

	mizuki.Core.Tasks.Add(func() {
		p.char.Core.Log.NewEvent("Picked up snack", glog.LogCharacterEvent, activeChar.Index).
			Write("heal", heal).
			Write("dmg", dmg)
		if dmg {
			p.explode()
		}

		if heal {
			// Heals double the amount on Mizuki
			healMultiplier := 1.0
			if activeChar.Index == mizuki.Index {
				healMultiplier = 2.0
			}
			mizuki.Core.Player.Heal(info.HealInfo{
				Caller:  mizuki.Index,
				Target:  activeChar.Index,
				Message: snackHealName,
				Src:     ((mizuki.Stat(attributes.EM) * snackHealEM[mizuki.TalentLvlBurst()]) + snackHealFlat[mizuki.TalentLvlBurst()]) * healMultiplier,
				Bonus:   mizuki.Stat(attributes.Heal),
			})
		}

		// C4 restores 5 energy to mizuki up to 4 times
		if mizuki.Base.Cons >= 4 {
			if mizuki.c4EnergyGenerationsRemaining > 0 {
				mizuki.c4EnergyGenerationsRemaining--
				mizuki.AddEnergy(c4Key, c4Energy)
			}
		}
	}, snackHitmark)
	p.Kill()
}

func (p *snack) explode() {
	p.explodeWithHitmark(0)
}

func (p *snack) explodeWithHitmark(hitmark int) {
	p.char.Core.QueueAttackWithSnap(p.attackInfo, p.snapshot, p.pattern, hitmark)
}

func (p *snack) HandleAttack(atk *combat.AttackEvent) float64 {
	// only collisions with the player can affect this or if it expires
	return 0
}

func (p *snack) SetDirection(trg geometry.Point) {}
func (p *snack) SetDirectionToClosestEnemy()     {}
func (p *snack) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}
