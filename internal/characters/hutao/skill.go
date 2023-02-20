package hutao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const (
	skillStart        = 14
	paramitaBuff      = "paramita"
	paramitaEnergyICD = "paramita-ball-icd"
	bbDebuff          = "blood-blossom"
)

func init() {
	skillFrames = frames.InitAbilSlice(52)
	skillFrames[action.ActionAttack] = 29
	skillFrames[action.ActionBurst] = 28
	skillFrames[action.ActionDash] = 37
	skillFrames[action.ActionJump] = 37
}

func (c *char) Skill(p map[string]int) action.ActionInfo {

	bonus := ppatk[c.TalentLvlSkill()] * c.MaxHP()
	max := (c.Base.Atk + c.Weapon.Atk) * 4
	if bonus > max {
		bonus = max
	}
	c.ppbuff[attributes.ATK] = bonus
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(paramitaBuff, 540+skillStart),
		AffectedStat: attributes.ATK,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			return c.ppbuff, true
		},
	})
	//TODO; this applies a1 at the end of paramita without checking for "pp extend" (if that's real)
	c.applyA1 = true
	c.QueueCharTask(c.a1, 540+skillStart)

	//remove some hp
	c.Core.Player.Drain(player.DrainInfo{
		ActorIndex: c.Index,
		Abil:       "Paramita Papilio",
		Amount:     .30 * c.HPCurrent,
	})

	//trigger 0 damage attack; matters because this breaks freeze
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Paramita (0 dmg)",
		AttackTag:  combat.AttackTagNone,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Physical,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3), skillStart, skillStart)

	c.SetCDWithDelay(action.ActionSkill, 960, 14)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if !c.StatModIsActive(paramitaBuff) {
		return
	}
	if c.StatusIsActive(paramitaEnergyICD) {
		return
	}
	c.AddStatus(paramitaEnergyICD, 5*60, true)

	count := 2.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 3
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Pyro, c.ParticleDelay) // TODO: this used to be 80
}

func (c *char) applyBB(a combat.AttackCB) {
	trg, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	if !trg.StatusIsActive(bbDebuff) {
		//start ticks
		trg.QueueEnemyTask(c.bbtickfunc(c.Core.F, trg), 240)
		trg.SetTag(bbDebuff, c.Core.F) //to track current bb source
	}

	trg.AddStatus(bbDebuff, 570, true) //lasts 8s + 1.5s
}

func (c *char) bbtickfunc(src int, trg *enemy.Enemy) func() {
	return func() {
		//do nothing if source changed
		if trg.Tags[bbDebuff] != src {
			return
		}
		if !trg.StatusIsActive(bbDebuff) {
			return
		}
		c.Core.Log.NewEvent("Blood Blossom checking for tick", glog.LogCharacterEvent, c.Index).
			Write("src", src)

		//queue up one damage instance
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Blood Blossom",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       bb[c.TalentLvlSkill()],
		}
		//if cons 2, add flat dmg
		if c.Base.Cons >= 2 {
			ai.FlatDmg += c.MaxHP() * 0.1
		}
		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(trg.Key()), 0, 0)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Blood Blossom ticked", glog.LogCharacterEvent, c.Index).
				Write("next expected tick", c.Core.F+240).
				Write("dur", trg.StatusExpiry(bbDebuff)).
				Write("src", src)
		}
		//queue up next instance
		trg.QueueEnemyTask(c.bbtickfunc(src, trg), 240)
	}
}
