package core

type WeaponClass int

const (
	WeaponClassSword WeaponClass = iota
	WeaponClassClaymore
	WeaponClassSpear
	WeaponClassBow
	WeaponClassCatalyst
	EndWeaponClass
)

var weaponName = []string{
	"sword",
	"claymore",
	"polearm",
	"bow",
	"catalyst",
}

func (w WeaponClass) String() string {
	return weaponName[w]
}
