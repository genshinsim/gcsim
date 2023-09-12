package glog

type Source int

const (
	LogInvalid     Source = iota
	LogDamageEvent        // tracks damages
	LogPreDamageMod
	LogCalc         // detailed calcs
	LogElementEvent // tracks elemental application
	LogSnapshotEvent
	LogStatusEvent
	LogActionEvent
	LogEnergyEvent
	LogCharacterEvent
	LogEnemyEvent
	LogSimEvent
	LogArtifactEvent
	LogWeaponEvent
	LogShieldEvent
	LogConstructEvent
	LogICDEvent
	LogDebugEvent    // for debug purposes
	LogWarnings      // for things that go wrong
	LogPlayerEvent   // for events related to player i.e. stamina, swap cd, healing, taking dmg etc..
	LogHealEvent     // healing events
	LogHurtEvent     // taking dmg event
	LogCooldownEvent // for tracking things going on and off cooldown
	LogHitlagEvent   // for debugging hitlag
	LogUserEvent     // user print event
)

var LogSourceFromString = map[string]Source{
	"":                LogInvalid,
	"damage":          LogDamageEvent,
	"pre_damage_mods": LogPreDamageMod,
	"calc":            LogCalc,
	"element":         LogElementEvent,
	"snapshot":        LogSnapshotEvent,
	"status":          LogStatusEvent,
	"action":          LogActionEvent,
	"energy":          LogEnergyEvent,
	"character":       LogCharacterEvent,
	"enemy":           LogEnemyEvent,
	"sim":             LogSimEvent,
	"artifact":        LogArtifactEvent,
	"weapon":          LogWeaponEvent,
	"shield":          LogShieldEvent,
	"construct":       LogConstructEvent,
	"icd":             LogICDEvent,
	"debug":           LogDebugEvent,
	"warning":         LogWarnings,
	"player":          LogPlayerEvent,
	"heal":            LogHealEvent,
	"hurt":            LogHurtEvent,
	"cooldown":        LogCooldownEvent,
	"hitlag":          LogHitlagEvent,
	"user":            LogUserEvent,
}

var LogSourceString = [...]string{
	"",
	"damage",
	"pre_damage_mods",
	"calc",
	"element",
	"snapshot",
	"status",
	"action",
	"energy",
	"character",
	"enemy",
	"sim",
	"artifact",
	"weapon",
	"shield",
	"construct",
	"icd",
	"debug",
	"warning",
	"player",
	"heal",
	"hurt",
	"cooldown",
	"hitlag",
	"user",
}

func (l Source) String() string {
	return LogSourceString[l]
}
