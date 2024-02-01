package attacks

type AimParam int

const (
	AimParamPhys = iota // physical aimed shot (no arcc)
	AimParamLv1         // fully charged aimshot for most bow characters
	AimParamLv2         // used by characters like ganyu or yoimiya
	AimParamLv3         // yoimiya only for now (2 kindling arrows)
	AimParamLv4         // yoimiya only for now (3 kindling arrows)
)
