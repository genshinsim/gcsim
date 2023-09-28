package conditional

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func evalGadgets(c *core.Core, fields []string) (int, error) {
	if err := fieldsCheck(fields, 3, "gadgets"); err != nil {
		return 0, err
	}
	switch fields[1] {
	case "dendrocore":
		return evalDendroCore(c, fields[2])
	default:
		return 0, fmt.Errorf("bad gadgets condition: invalid criteria %v", fields[1])
	}
}

func evalDendroCore(c *core.Core, key string) (int, error) {
	switch key {
	case "count":
		count := 0
		for i := 0; i < c.Combat.GadgetCount(); i++ {
			if _, ok := c.Combat.Gadget(i).(*reactable.DendroCore); ok {
				count++
			}
		}
		return count, nil
	default:
		return 0, fmt.Errorf("bad gadgets (dendrocore) condition: invalid criteria %v", key)
	}
}
