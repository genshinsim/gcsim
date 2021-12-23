package travelergeo

import "github.com/genshinsim/gcsim/pkg/core"

type stone struct {
	src    int
	expiry int
	char   *char
}

func (s *stone) Key() int {
	return s.src
}

func (s *stone) Type() core.GeoConstructType {
	return core.GeoConstructTravellerSkill
}

func (s *stone) OnDestruct() {
	if s.char.Base.Cons >= 2 {
		ai := core.AttackInfo{
			ActorIndex: s.char.Index,
			Abil:       "Rockcore Meltdown",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeBlunt,
			Element:    core.Geo,
			Durability: 50,
			Mult:       skill[s.char.TalentLvlSkill()],
		}
		s.char.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 0)
	}
}

func (s *stone) Expiry() int {
	return s.expiry
}

func (s *stone) IsLimited() bool {
	return true
}

func (s *stone) Count() int {
	return 1
}

type barrier struct {
	src    int
	expiry int
	char   *char
}

func (b *barrier) Key() int {
	return b.src
}

func (b *barrier) Type() core.GeoConstructType {
	return core.GeoConstructTravellerBurst
}

func (b *barrier) OnDestruct() {
	if b.char.Base.Cons >= 1 {
		b.char.Tags["wall"] = 0
	}
}

func (b *barrier) Expiry() int {
	return b.expiry
}

func (b *barrier) IsLimited() bool {
	return true
}

func (b *barrier) Count() int {
	return 3
}
