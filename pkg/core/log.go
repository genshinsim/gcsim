package core

type LogEvent interface {
	LogSource() LogSource           //returns the type of this log event i.e. character, sim, damage, etc...
	StartFrame() int                //returns the frame on which this event was started
	Src() int                       //returns the index of the character that triggered this event. -1 if it's not a character
	Write(keyAndVal ...interface{}) //write additional keyAndVal pairs to the event
	SetEnded(f int)
}

type LogCtrl interface {
	NewEvent(msg string, typ LogSource, srcChar int, keysAndValues ...interface{}) LogEvent
	NewEventBuildMsg(typ LogSource, srcChar int, msg ...string) LogEvent
	Dump() ([]byte, error) //print out all the logged events in array of JSON strings in the ordered they were added
}

type NilLogEvent struct{}

func (n *NilLogEvent) LogSource() LogSource           { return LogSimEvent }
func (n *NilLogEvent) StartFrame() int                { return -1 }
func (n *NilLogEvent) Src() int                       { return 0 }
func (n *NilLogEvent) Write(keyAndVal ...interface{}) {}
func (n *NilLogEvent) SetEnded(f int)                 {}

type NilLogger struct{}

func (n *NilLogger) Dump() ([]byte, error) { return []byte{}, nil }
func (n *NilLogger) NewEventBuildMsg(typ LogSource, srcChar int, msg ...string) LogEvent {
	return &NilLogEvent{}
}
func (n *NilLogger) NewEvent(msg string, typ LogSource, srcChar int, keysAndValues ...interface{}) LogEvent {
	return &NilLogEvent{}
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
