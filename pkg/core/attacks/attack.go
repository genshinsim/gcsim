package attacks

import "slices"

type AttackTag int // attacktag is used instead of actions etc..

const (
	AttackTagNone AttackTag = iota
	AttackTagNormal
	AttackTagExtra
	AttackTagPlunge
	AttackTagElementalArt
	AttackTagElementalArtHold
	AttackTagElementalBurst
	AttackTagRelicSkill
	AttackTagWeaponSkill
	AttackTagMonaBubbleBreak

	AttackTagNoneStat // ignore attacker stats delim

	AttackTagOverloadDamage
	AttackTagSuperconductDamage
	AttackTagECDamage
	AttackTagShatter
	AttackTagSwirlPyro
	AttackTagSwirlHydro
	AttackTagSwirlCryo
	AttackTagSwirlElectro
	AttackTagBurningDamage
	AttackTagBloom
	AttackTagBountifulCore // special tag for nilou
	AttackTagBurgeon
	AttackTagHyperbloom

	AttackTagReactionLunarCharge
	AttackTagReactionLunarCrystallize

	AttackTagDirectLunarCharged
	AttackTagDirectLunarBloom
	AttackTagDirectLunarCrystallize

	AttackTagReactionStellarSwirl

	AttackTagDirectStellarConduct
	AttackTagDirectStellarSwirl
)

// TODO: get rid of direct/reaction split
func (t AttackTag) IsReaction() bool { return t >= AttackTagNoneStat }
func (t AttackTag) IsDirect() bool   { return t.IsLunarDirect() || t.IsStellarDirect() }

var (
	lunarReact = []AttackTag{
		AttackTagReactionLunarCharge,
		AttackTagReactionLunarCrystallize,
	}
	lunarDirect = []AttackTag{
		AttackTagDirectLunarCharged,
		AttackTagDirectLunarBloom,
		AttackTagDirectLunarCrystallize,
	}
)

func (t AttackTag) IsLunar() bool       { return t.IsLunarReact() || t.IsLunarDirect() }
func (t AttackTag) IsLunarReact() bool  { return slices.Contains(lunarReact, t) }
func (t AttackTag) IsLunarDirect() bool { return slices.Contains(lunarDirect, t) }

var (
	stellarReact = []AttackTag{
		AttackTagReactionStellarSwirl,
	}
	stellarDirect = []AttackTag{
		AttackTagDirectStellarConduct,
		AttackTagDirectStellarSwirl,
	}
)

func (t AttackTag) IsStellar() bool       { return t.IsStellarReact() || t.IsStellarDirect() }
func (t AttackTag) IsStellarReact() bool  { return slices.Contains(stellarReact, t) }
func (t AttackTag) IsStellarDirect() bool { return slices.Contains(stellarDirect, t) }

type StrikeType int

const (
	StrikeTypeDefault StrikeType = iota
	StrikeTypePierce
	StrikeTypeBlunt
	StrikeTypeSlash
	StrikeTypeSpear
)

// TODO: merge AdditionalTag into AttackTag
type AdditionalTag int

const (
	AdditionalTagNone AdditionalTag = iota
	AdditionalTagNightsoul
	AdditionalTagKinichCannon
	AdditionalTagVarkaSpecial
)
