package keys

import (
	"encoding/json"
	"fmt"
)

type TargetID int

type TargetIDGenerator struct {
	state int
}

func NewTargetIDGenerator() *TargetIDGenerator {
	return &TargetIDGenerator{
		state: 1, // start at 1 to ensure 0 is always invalid
	}
}

// Generate a new TargetID
func (gen *TargetIDGenerator) New() TargetID {
	out := gen.state
	gen.state += 1
	return TargetID(out)
}

// Size returns the number of ids handed out, for purpose of
// sizing arrays etc...
func (gen *TargetIDGenerator) Size() int {
	return gen.state - 1 // because we start at 0
}

func (t TargetID) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprint(t))
}
