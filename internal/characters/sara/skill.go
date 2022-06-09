package sara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillFrames []int

// c2 hitmark
const skillHitmark = 90

// Implements skill handling. Fairly barebones since most of the actual stuff happens elsewhere
// Gains Crowfeather Cover for 18s, and when Kujou Sara fires a fully-charged Aimed Shot, Crowfeather Cover will be consumed, and will leave a Crowfeather at the target location.
// Crowfeathers will trigger Tengu Juurai: Ambush after a short time, dealing Electro DMG and granting the active character within its AoE an ATK Bonus based on Kujou Sara's Base ATK.
// The ATK Bonuses from different Tengu Juurai will not stack, and their effects and duration will be determined by the last Tengu Juurai to take effect.
// Also implements C2: Unleashing Tengu Stormcall will leave a Weaker Crowfeather at Kujou Sara's original position that will deal 30% of its original DMG.
func (c *char) Skill(p map[string]int) action.ActionInfo {

	// Snapshot for all of the crowfeathers are taken upon cast
	c.Core.Status.Add("saracover", 18*60)

	// C2 handling
	if c.Base.Cons >= 2 {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Tengu Juurai: Ambush C2",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.3 * skill[c.TalentLvlSkill()],
		}
		// TODO: not sure of snapshot? timing
		c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 50, skillHitmark, c.a4)
		c.attackBuff(skillHitmark)
	}

	c.SetCD(action.ActionSkill, 600)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		Post:            skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

// Handles attack boost from Sara's skills
// Checks for the onfield character at the delay frame, then applies buff to that character
func (c *char) attackBuff(delay int) {
	c.Core.Tasks.Add(func() {
		buff := atkBuff[c.TalentLvlSkill()] * float64(c.Base.Atk+c.Weapon.Atk)

		active := c.Core.Player.ActiveChar()
		active.SetTag("sarabuff", c.Core.F+360)
		c.Core.Log.NewEvent("sara attack buff applied", glog.LogCharacterEvent, c.Index, "char", active.Index, "buff", buff, "expiry", c.Core.F+360)

		m := make([]float64, attributes.EndStatType)
		m[attributes.ATK] = buff
		active.AddStatMod("sara-attack-buff", 360, attributes.ATK, func() ([]float64, bool) {
			return m, true
		})

		if c.Base.Cons >= 1 {
			c.c1()
		}
		if c.Base.Cons >= 6 {
			c.c6(active)
		}
	}, delay)
}
