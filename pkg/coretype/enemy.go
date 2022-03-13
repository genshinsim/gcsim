package coretype

type Enemy interface {
	Indexer //index for tracking enemy number
	Tagger  //tags
	Target
	Reactable

	//target mods
	AddDefMod(key string, val float64, dur int)
	AddResMod(key string, val ResistMod)
	RemoveResMod(key string)
	RemoveDefMod(key string)
	HasDefMod(key string) bool
	HasResMod(key string) bool
}

type ResistMod struct {
	Key      string
	Ele      EleType
	Value    float64
	Duration int
	Expiry   int
}
