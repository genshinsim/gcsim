package targets

type TargetKey int

const InvalidTargetKey TargetKey = -1

type TargettableType int

const (
	TargettableEnemy TargettableType = iota
	TargettablePlayer
	TargettableGadget
	TargettableTypeCount
)
