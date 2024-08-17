package info

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Sets map[keys.Set]int

func (s Sets) MarshalJSON() ([]byte, error) {
	// we'll use a custom string builder i guess
	var sb strings.Builder
	sb.WriteString("{")
	for k, v := range s {
		sb.WriteString(`"`)
		sb.WriteString(k.String())
		sb.WriteString(`":"`)
		sb.WriteString(strconv.Itoa(v))
		sb.WriteString(`",`)
	}
	str := sb.String()
	str = strings.TrimRight(str, ",")
	str += "}"
	return []byte(str), nil
}

type Set interface {
	SetIndex(int)
	GetCount() int
	Init() error
}
