package keys

type Set int

func (c Set) String() string {
	return setNames[c]
}

var setNames = []string{}
