package coretype

type LogEvent interface {
	LogSource() LogSource           //returns the type of this log event i.e. character, sim, damage, etc...
	StartFrame() int                //returns the frame on which this event was started
	Src() int                       //returns the index of the character that triggered this event. -1 if it's not a character
	Write(keyAndVal ...interface{}) //write additional keyAndVal pairs to the event
	SetEnded(f int)
}

type LogSource int

const (
	LogProcs LogSource = iota
	LogDamageEvent
	LogPreDamageMod
	LogHurtEvent
	LogHealEvent
	LogCalc
	LogReactionEvent
	LogElementEvent
	LogSnapshotEvent
	LogSnapshotModsEvent
	LogStatusEvent
	LogActionEvent
	LogQueueEvent
	LogEnergyEvent
	LogCharacterEvent
	LogEnemyEvent
	LogHookEvent
	LogSimEvent
	LogTaskEvent
	LogArtifactEvent
	LogWeaponEvent
	LogShieldEvent
	LogConstructEvent
	LogICDEvent
)

var LogSourceFromString = map[string]LogSource{
	"procs":           LogProcs,
	"damage":          LogDamageEvent,
	"pre_damage_mods": LogPreDamageMod,
	"hurt":            LogHurtEvent,
	"heal":            LogHealEvent,
	"calc":            LogCalc,
	"reaction":        LogReactionEvent,
	"element":         LogElementEvent,
	"snapshot":        LogSnapshotEvent,
	"snapshot_mods":   LogSnapshotModsEvent,
	"status":          LogStatusEvent,
	"action":          LogActionEvent,
	"queue":           LogQueueEvent,
	"energy":          LogEnergyEvent,
	"character":       LogCharacterEvent,
	"enemy":           LogEnemyEvent,
	"hook":            LogHookEvent,
	"sim":             LogSimEvent,
	"task":            LogTaskEvent,
	"artifact":        LogArtifactEvent,
	"weapon":          LogWeaponEvent,
	"shield":          LogShieldEvent,
	"construct":       LogConstructEvent,
	"icd":             LogICDEvent,
}

var LogSourceString = [...]string{
	"procs",
	"damage",
	"pre_damage_mods",
	"hurt",
	"heal",
	"calc",
	"reaction",
	"element",
	"snapshot",
	"snapshot_mods",
	"status",
	"action",
	"queue",
	"energy",
	"character",
	"enemy",
	"hook",
	"sim",
	"task",
	"artifact",
	"weapon",
	"shield",
	"construct",
	"icd",
}

func (l LogSource) String() string {
	return LogSourceString[l]
}
