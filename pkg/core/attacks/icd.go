package attacks

type ICDTag int // same ICD tag shares the same counter

const (
	ICDTagNone ICDTag = iota
	ICDTagNormalAttack
	ICDTagExtraAttack
	ICDTagElementalArt
	ICDTagElementalArtHold
	ICDTagElementalArtAnemo
	ICDTagElementalArtPyro
	ICDTagElementalArtHydro
	ICDTagElementalArtCryo
	ICDTagElementalArtElectro
	ICDTagElementalBurst
	ICDTagElementalBurstAnemo
	ICDTagElementalBurstPyro
	ICDTagElementalBurstHydro
	ICDTagElementalBurstCryo
	ICDTagElementalBurstElectro
	ICDTagDash
	ICDTagTravelerWakeOfEarth
	ICDTagTravelerDewdrop
	ICDTagTravelerScorchingThreshold
	ICDTagTravelerBlazingThreshold
	ICDTagTravelerHoldDMG

	ICDReactionDamageDelim
	ICDTagOverloadDamage
	ICDTagSuperconductDamage
	ICDTagECDamage
	ICDTagShatter
	ICDTagSwirlPyro
	ICDTagSwirlHydro
	ICDTagSwirlCryo
	ICDTagSwirlElectro
	ICDTagBurningDamage
	ICDTagBloomDamage
	ICDTagBurgeonDamage
	ICDTagHyperbloomDamage

	EndDefaultICDTags
)

// group dictate both the sequence and the reset timer
type ICDGroup int // same ICD group shares the same timer

const (
	ICDGroupDefault ICDGroup = iota
	ICDGroupPoleExtraAttack
	ICDGroupReactionA
	ICDGroupReactionB
	ICDGroupBurning
	ICDGroupTravelerDewdrop
	ICDGroupTravelerBurst
	EndDefaultICDGroups
)

var ICDGroupResetTimer []int
var ICDGroupEleApplicationSequence [][]float64
var ICDGroupDamageSequence [][]float64

func init() {
	ICDGroupResetTimer = make([]int, ICDGroupLength)
	ICDGroupResetTimer[ICDGroupDefault] = 150
	ICDGroupResetTimer[ICDGroupPoleExtraAttack] = 30
	ICDGroupResetTimer[ICDGroupReactionA] = 30
	ICDGroupResetTimer[ICDGroupReactionB] = 30
	ICDGroupResetTimer[ICDGroupBurning] = 120
	ICDGroupResetTimer[ICDGroupTravelerDewdrop] = 90
	ICDGroupResetTimer[ICDGroupTravelerBurst] = 480

	ICDGroupEleApplicationSequence = make([][]float64, ICDGroupLength)
	ICDGroupEleApplicationSequence[ICDGroupDefault] = []float64{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupPoleExtraAttack] = []float64{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupReactionA] = []float64{1, 1}
	ICDGroupEleApplicationSequence[ICDGroupReactionB] = []float64{1, 1}
	ICDGroupEleApplicationSequence[ICDGroupBurning] = []float64{1, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupTravelerDewdrop] = []float64{1, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupTravelerBurst] = []float64{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}

	ICDGroupDamageSequence = make([][]float64, ICDGroupLength)
	ICDGroupDamageSequence[ICDGroupDefault] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupPoleExtraAttack] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupDamageSequence[ICDGroupReactionA] = []float64{1, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupDamageSequence[ICDGroupReactionB] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	// actual data: {1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0}
	// however there seems to be no limit to the amount of burning dmg a target can take
	ICDGroupDamageSequence[ICDGroupBurning] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupTravelerDewdrop] = []float64{1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupTravelerBurst] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
}
