package keys

type Weapon int

func (c Weapon) String() string {
	return weaponNames[c]
}

var weaponNames = []string{}
