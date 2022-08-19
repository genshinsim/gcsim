package conditional

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/construct"
)

func evalConstruct[V core.Number](c *core.Core, fields []string) (V, error) {
	if err := fieldsCheck(fields, 3, "construct"); err != nil {
		return 0, err
	}
	switch fields[1] {
	case ".duration":
		return evalConstructDuration[V](c, fields)
	case ".count":
		return evalConstructCount[V](c, fields)
	default:
		return 0, fmt.Errorf("bad construct condition: invalid criteria %v", fields[1])
	}
}

func evalConstructDuration[V core.Number](c *core.Core, fields []string) (V, error) {
	//.construct.duration.<name>
	s := strings.TrimPrefix(fields[2], ".")
	key, ok := construct.ConstructNameToKey[s]
	if !ok {
		return 0, fmt.Errorf("bad construction condition: invalid construct %v", s)
	}
	return V(c.Constructs.Expiry(key)), nil
}

func evalConstructCount[V core.Number](c *core.Core, fields []string) (V, error) {
	//.construct.count.<name>
	s := strings.TrimPrefix(fields[2], ".")
	key, ok := construct.ConstructNameToKey[s]
	if !ok {
		return 0, fmt.Errorf("bad construction condition: invalid construct %v", s)
	}
	return V(c.Constructs.CountByType(key)), nil
}
