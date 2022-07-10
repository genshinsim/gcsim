package diluc

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstHitmark = 100

func init() {
	burstFrames = frames.InitAbilSlice(141)
	burstFrames[action.ActionAttack] = 140
	burstFrames[action.ActionSkill] = 139
	burstFrames[action.ActionDash] = 139
	burstFrames[action.ActionSwap] = 138
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	dot, ok := p["dot"]
	if !ok {
		dot = 2 //number of dot hits
	}
	if dot > 7 {
		dot = 7
	}
	explode, ok := p["explode"]
	if !ok {
		explode = 0 //if explode hits
	}

	//enhance weapon for 12 seconds (with a4)
	// Infusion starts when burst starts and ends when burst comes off CD - check any diluc video
	c.Core.Status.Add("dilucq", 720)
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		"diluc-fire-weapon",
		attributes.Pyro,
		720,
		false,
		combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
	)

	// a4: add 20% pyro damage
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("diluc-fire-weapon", 720),
		AffectedStat: attributes.PyroP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	// Snapshot occurs late in the animation when it is released from the claymore
	// For our purposes, snapshot upon damage proc
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Dawn (Strike)",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagElementalBurst,
			ICDGroup:   combat.ICDGroupDiluc,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Pyro,
			Durability: 50,
			Mult:       burstInitial[c.TalentLvlBurst()],
		}

		c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, 1)

		//dot does damage every .2 seconds for 7 hits? so every 12 frames
		//dot does max 7 hits + explosion, roughly every 13 frame? blows up at 210 frames
		//first tick did 50 dur as well?
		ai.Abil = "Dawn (Tick)"
		ai.Mult = burstDOT[c.TalentLvlBurst()]
		for i := 1; i <= dot; i++ {
			c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, i+12)
		}

		if explode > 0 {
			ai.Abil = "Dawn (Explode)"
			ai.Mult = burstExplode[c.TalentLvlBurst()]
			c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, 110)
		}
	}, burstHitmark)

	c.ConsumeEnergy(21)
	c.SetCDWithDelay(action.ActionBurst, 720, 14)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}
