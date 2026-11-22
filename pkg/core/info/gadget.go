package info

type GadgetTyp int

const (
	GadgetTypUnknown GadgetTyp = iota
	StartGadgetTypEnemy
	GadgetTypDendroCore
	GadgetTypLeaLotus
	GadgetTypBogglecatBox
	EndGadgetTypEnemy
	GadgetTypGuoba
	GadgetTypYueguiThrowing
	GadgetTypYueguiJumping
	GadgetTypBaronBunny
	GadgetTypGrinMalkinHat
	GadgetTypSourcewaterDropletHydroTrav
	GadgetTypSourcewaterDropletNeuv
	GadgetTypSourcewaterDropletSigewinne
	GadgetTypCrystallizeShard
	GadgetTypYumemiSnack
	GadgetTypPolestarField
	GadgetTypStellarVortex
	GadgetTypTest
	EndGadgetTyp
)

type Gadget interface {
	Target
	Src() int
	GadgetTyp() GadgetTyp
	// HandleSharedAttack takes an attack like HandleAttack but without a copy of its own of the
	// event, and returns false, without doing anything, if it needs one. It may read a but has
	// to copy it before changing it or passing it on to anything that could keep it.
	HandleSharedAttack(a *AttackEvent) (float64, bool)
}
