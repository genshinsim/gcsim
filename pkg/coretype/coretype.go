package coretype

type Flags struct {
	DamageMode     bool
	EnergyCalcMode bool // Allows Burst Action when not at full Energy, logs current Energy when using Burst
	LogDebug       bool // Used to determine logging level
	ChildeActive   bool // Used for Childe +1 NA talent passive
	Delays         Delays
	Custom         map[string]int
}

type Delays struct {
	Skill  int
	Burst  int
	Attack int
	Charge int
	Aim    int
	Dash   int
	Jump   int
	Swap   int
}
