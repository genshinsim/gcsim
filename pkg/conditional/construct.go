package conditional

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/construct"
)

func evalConstruct(c *core.Core, fields []string) (int, error) {
	// .construct.count.<name>
	// .construct.duration.<name>
	if err := fieldsCheck(fields, 3, "construct"); err != nil {
		return 0, err
	}

	name := fields[2]
	key, ok := construct.ConstructNameToKey[name]
	if !ok {
		return 0, fmt.Errorf("bad construct condition: invalid construct %v", name)
	}

	switch v := fields[1]; v {
	case countField:
		return c.Constructs.CountByType(key), nil
	case "duration":
		return c.Constructs.Expiry(key), nil
	default:
		return 0, fmt.Errorf("bad construct condition: invalid criteria %v", v)
	}
}
