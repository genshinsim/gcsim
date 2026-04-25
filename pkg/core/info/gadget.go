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
	GadgetTypTest
	EndGadgetTyp
)

type Gadget interface {
	Target
	Src() int
	GadgetTyp() GadgetTyp
}
