package def

type Target interface {
	SetIndex(ind int) //update the current index
	MaxHP() float64
	HP() float64
	Tick()
	Attack(ds *Snapshot) float64

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
