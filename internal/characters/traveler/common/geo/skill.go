package geo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames [][][]int
	// {Tap E, Short Hold E}
	skillHitmark = []int{62, 29}
	skillCDStart = []int{23, 25}
)

func init() {
	skillFrames = make([][][]int, 2)

	// Tap E
	skillFrames[0] = make([][]int, 2)

	// Male
	skillFrames[0][0] = frames.InitAbilSlice(81) // Tap E -> N1/Q
	skillFrames[0][0][action.ActionDash] = 25    // Tap E -> D
	skillFrames[0][0][action.ActionJump] = 24    // Tap E -> J
	skillFrames[0][0][action.ActionSwap] = 67    // Tap E -> Swap

	// Female
	skillFrames[0][1] = frames.InitAbilSlice(80) // Tap E -> Q
	skillFrames[0][1][action.ActionAttack] = 79  // Tap E -> N1
	skillFrames[0][1][action.ActionDash] = 23    // Tap E -> D
	skillFrames[0][1][action.ActionJump] = 24    // Tap E -> J
	skillFrames[0][1][action.ActionSwap] = 65    // Tap E -> Swap

	// Short Hold E
	skillFrames[1] = make([][]int, 2)

	// Male
	skillFrames[1][0] = frames.InitAbilSlice(54) // Short Hold E -> N1/Q
	skillFrames[1][0][action.ActionDash] = 31    // Short Hold E -> D
	skillFrames[1][0][action.ActionJump] = 31    // Short Hold E -> J
	skillFrames[1][0][action.ActionSwap] = 39    // Short Hold E -> Swap

	// Female
	skillFrames[1][1] = frames.InitAbilSlice(54) // Short Hold E -> N1/Q
	skillFrames[1][1][action.ActionDash] = 31    // Short Hold E -> D
	skillFrames[1][1][action.ActionJump] = 32    // Short Hold E -> J
	skillFrames[1][1][action.ActionSwap] = 40    // Short Hold E -> Swap
}

func (c *Traveler) Skill(p map[string]int) (action.Info, error) {
	shortHold, ok := p["short_hold"]
	if !ok || shortHold < 0 {
		shortHold = 0
	}
	if shortHold > 1 {
		shortHold = 1
	}

	noMeteorite := p["no_meteorite"] == 1

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Starfell Sword",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           133.2,
		Element:            attributes.Geo,
		Durability:         50,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagFactor:       0.05,
		HitlagHaltFrames:   0.05 * 60,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	stoneDir := c.Core.Combat.Player().Direction()
	stonePos := c.Core.Combat.PrimaryTarget().Pos()

	// TODO: check snapshot timing
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(stonePos, nil, 3),
		24,
		skillHitmark[shortHold],
		c.makeParticleCB(),
	)

	if !noMeteorite {
		c.Core.Tasks.Add(func() {
			dur := 30 * 60
			if c.Base.Cons >= 6 {
				dur += 600
			}
			c.Core.Constructs.New(c.newStone(dur, stoneDir, stonePos), false)
		}, skillHitmark[shortHold])
	}

	c.SetCDWithDelay(action.ActionSkill, c.skillCD, skillCDStart[shortHold])

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[shortHold][c.gender]),
		AnimationLength: skillFrames[shortHold][c.gender][action.InvalidAction],
		CanQueueAfter:   skillFrames[shortHold][c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *Traveler) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		count := 3.0
		if c.Core.Rand.Float64() < 0.33 {
			count = 4
		}
		c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Geo, c.ParticleDelay)
	}
}

type stone struct {
	src    int
	expiry int
	char   *Traveler
	dir    geometry.Point
	pos    geometry.Point
}

func (c *Traveler) newStone(dur int, dir, pos geometry.Point) *stone {
	return &stone{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
		dir:    dir,
		pos:    pos,
	}
}

func (s *stone) OnDestruct() {
	if s.char.Base.Cons >= 2 {
		ai := combat.AttackInfo{
			ActorIndex:         s.char.Index,
			Abil:               "Rockcore Meltdown",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagElementalArt,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           133.2,
			Element:            attributes.Geo,
			Durability:         50,
			Mult:               skill[s.char.TalentLvlSkill()],
			HitlagFactor:       0.05,
			HitlagHaltFrames:   0.05 * 60,
			CanBeDefenseHalted: true,
			IsDeployable:       true,
		}
		s.char.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(s.pos, nil, 3),
			0,
			0,
		)
	}
}

func (s *stone) Key() int                         { return s.src }
func (s *stone) Type() construct.GeoConstructType { return construct.GeoConstructTravellerSkill }
func (s *stone) Expiry() int                      { return s.expiry }
func (s *stone) IsLimited() bool                  { return true }
func (s *stone) Count() int                       { return 1 }
func (s *stone) Direction() geometry.Point        { return s.dir }
func (s *stone) Pos() geometry.Point              { return s.pos }
