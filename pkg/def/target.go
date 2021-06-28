package def

type Target interface {
	SetIndex(ind int) //update the current index
	MaxHP() float64
	HP() float64
	AuraTick() //tick this first to avoid messing with combat
	Tick()
	Attack(ds *Snapshot) float64
	AddOnAttackLandedHook(fun func(ds *Snapshot), key string)
	RemoveOnAttackLandedHook(key string)

	Delete() //gracefully deference everything so that it can be gc'd
}

type ResistMod struct {
	Key      string
	Ele      EleType
	Value    float64
	Duration int
	Expiry   int
}

type DefMod struct {
	Key    int
	Value  float64
	Expiry int
}

type ReactionBonusMod struct {
	Key    string
	Amount func(ds *Snapshot) float64
}
