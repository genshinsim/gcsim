package info

import (
	"encoding/json"
	"fmt"
)

type EntityIndex int

type EntityIndexRegistry struct {
	lookup []int
}

func (gen *EntityIndexRegistry) Find(id int) EntityIndex {
	for i, v := range gen.lookup {
		if v == id {
			return EntityIndex(i + 1)
		}
	}
	return 0
}

// Generate a new TargetID
func (gen *EntityIndexRegistry) Register(id int) (EntityIndex, error) {
	if exist := gen.Find(id); exist != 0 {
		return 0, fmt.Errorf("id %v already exist", id)
	}
	gen.lookup = append(gen.lookup, id)
	return EntityIndex(len(gen.lookup)), nil // will always return index starting at 1
}

// Size returns the number of ids handed out, for purpose of
// sizing arrays etc...
func (gen *EntityIndexRegistry) Size() int {
	return len(gen.lookup) // because we start at 0
}

func (t EntityIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprint(t))
}
