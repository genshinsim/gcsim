package noelle

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = []int{28, 25, 20, 42}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	r := 0.3
	if c.Core.Status.Duration("noelleq") > 0 {
		r = 2
	}
	done := false
	cb := func(a core.AttackCB) {
		if done {
			return
		}
		//check for healing
		if c.Core.Shields.Get(core.ShieldNoelleSkill) != nil {
			var prob float64
			if c.Base.Cons >= 1 && c.Core.Status.Duration("noelleq") > 0 {
				prob = 1
			} else {
				prob = healChance[c.TalentLvlSkill()]
			}
			if c.Core.Rand.Float64() < prob {
				//heal target
				x := a.AttackEvent.Snapshot.BaseDef*(1+a.AttackEvent.Snapshot.Stats[core.DEFP]) + a.AttackEvent.Snapshot.Stats[core.DEF]
				heal := shieldHeal[c.TalentLvlSkill()]*x + shieldHealFlat[c.TalentLvlSkill()]
				c.Core.Health.Heal(core.HealInfo{
					Caller:  c.Index,
					Target:  -1,
					Message: "Breastplate (Attack)",
					Src:     heal,
					Bonus:   a.AttackEvent.Snapshot.Stats[core.Heal],
				})
				done = true
			}
		}

	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(r, false, core.TargettableEnemy), hitmarks[c.NormalCounter], hitmarks[c.NormalCounter], cb)

	c.AdvanceNormalIndex()

	c.a4Counter++
	if c.a4Counter == 4 {
		c.a4Counter = 0
		if c.Cooldown(core.ActionSkill) > 0 {
			c.ReduceActionCooldown(core.ActionSkill, 60)
		}
	}

	return f, a
}

type noelleShield struct {
	*shield.Tmpl
	c *char
}

func (n *noelleShield) OnExpire() {
	if n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
}

func (n *noelleShield) OnDamage(dmg float64, ele core.EleType, bonus float64) (float64, bool) {
	taken, ok := n.Tmpl.OnDamage(dmg, ele, bonus)
	if !ok && n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
	return taken, ok
}

func (c *char) newShield(base float64, t core.ShieldType, dur int) *noelleShield {
	n := &noelleShield{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.Src = c.Core.F
	n.Tmpl.ShieldType = t
	n.Tmpl.Name = "Noelle Skill"
	n.Tmpl.HP = base
	n.Tmpl.Expires = c.Core.F + dur
	n.c = c
	return n
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Breastplate",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 50,
		Mult:       shieldDmg[c.TalentLvlSkill()],
		UseDef:     true,
	}
	snap := c.Snapshot(&ai)

	//add shield first
	defFactor := snap.BaseDef*(1+snap.Stats[core.DEFP]) + snap.Stats[core.DEF]
	shield := shieldFlat[c.TalentLvlSkill()] + shieldDef[c.TalentLvlSkill()]*defFactor

	c.Core.Shields.Add(c.newShield(shield, core.ShieldNoelleSkill, 720))

	//activate shield timer, on expiry explode
	c.shieldTimer = c.Core.F + 720 //12 seconds

	c.a4Counter = 0

	x, y := c.Core.Targets[0].Shape().Pos()
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 2, false, core.TargettableEnemy), 15, 15)

	if c.Base.Cons >= 4 {
		c.AddTask(func() {
			if c.shieldTimer == c.Core.F {
				//deal damage
				c.explodeShield()
			}
		}, "noelle shield", 720)
	}

	c.SetCDWithDelay(core.ActionSkill, 24*60, 6)
	return f, a
}

func (c *char) explodeShield() {
	c.shieldTimer = 0
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Breastplate",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 50,
		Mult:       4,
	}

	x, y := c.Core.Targets[0].Shape().Pos()
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 4, false, core.TargettableEnemy), 0, 0)

}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	// Add mod for def to attack burst conversion
	// TODO: Assume snapshot happens immediately upon cast since the conversion buffs the two burst hits
	m := make([]float64, core.EndStatType)

	// Generate a "fake" snapshot in order to show a listing of the applied mods in the debug
	aiSnapshot := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sweeping Time (Stat Snapshot)",
	}
	snapshot := c.Snapshot(&aiSnapshot)
	burstDefSnapshot := snapshot.BaseDef*(1+snapshot.Stats[core.DEFP]) + snapshot.Stats[core.DEF]
	mult := defconv[c.TalentLvlBurst()]
	if c.Base.Cons >= 6 {
		mult += 0.5
	}
	m[core.ATK] = mult * burstDefSnapshot

	// TODO: Confirm exact timing of buff - for now matched to status duration previously set, which is 900 + animation frames
	c.AddMod(core.CharStatMod{
		Key:    "noelle-burst",
		Expiry: c.Core.F + 900 + f,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
	c.Core.Log.NewEvent("noelle burst", core.LogSnapshotEvent, c.Index, "total def", burstDefSnapshot, "atk added", m[core.ATK], "mult", mult)

	c.Core.Status.AddStatus("noelleq", 900+f)

	// Queue up task for Noelle burst extension
	// https://library.keqingmains.com/evidence/characters/geo/noelle#noelle-c6-burst-extension
	if c.Base.Cons >= 6 {
		c.AddTask(func() {
			if c.Core.ActiveChar == c.Index {
				return
			}
			c.Core.Log.NewEvent("noelle max burst extension activated", core.LogCharacterEvent, c.Index, "new_expiry", c.Core.F+600)

			// Adding the mod again with the same key replaces it
			c.AddMod(core.CharStatMod{
				Key:    "noelle-burst",
				Expiry: c.Core.F + 600,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
			c.Core.Status.AddStatus("noelleq", 600)
		}, "noelle-c6-burst-extension", 900+f)
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sweeping Time (Burst)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(6.5, false, core.TargettableEnemy), 24, 24)

	ai.Abil = "Sweeping Time (Skill)"
	ai.Mult = burstskill[c.TalentLvlBurst()]
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(4.5, false, core.TargettableEnemy), 65, 65)

	c.SetCD(core.ActionBurst, 900)
	c.ConsumeEnergy(8)
	return f, a
}
